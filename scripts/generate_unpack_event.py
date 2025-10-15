#!/usr/bin/env python3
"""
Go事件解析器生成工具

从Go文件中提取所有UnpackXXXEvent方法，生成统一的UnpackEvent方法、Topic0方法、MethodID方法和结构体Topic0方法，并直接写入到原文件。
"""

import re
import sys
import subprocess


def extract_unpack_methods(go_file_path):
    """提取所有UnpackXXXEvent方法信息"""
    try:
        with open(go_file_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except FileNotFoundError:
        print(f"错误：找不到文件 {go_file_path}")
        sys.exit(1)
    
    # 匹配UnpackXXXEvent方法的模式，同时提取Solidity事件签名（通用化结构体前缀）
    pattern = r'// Solidity: event ([^\n]+)\nfunc \((\w+) \*(\w+)\) (Unpack\w+Event)\(log \*types\.Log\) \(\*(\w+), error\)\s*{\s*event := "([^"]+)"'
    
    methods = []
    struct_prefix = None  # 用于存储结构体前缀（如T2, T3等）
    
    for match in re.finditer(pattern, content, re.DOTALL):
        solidity_signature = match.group(1).strip()
        receiver_var = match.group(2)  # 例如: t2, t3
        struct_name = match.group(3)   # 例如: T2, T3
        method_name = match.group(4)   # 例如: UnpackAuthorizedSignerUpdatedEvent
        return_type = match.group(5)   # 例如: T2AuthorizedSignerUpdated
        event_name = match.group(6)    # 例如: AuthorizedSignerUpdated
        
        # 动态确定结构体前缀
        if struct_prefix is None:
            struct_prefix = struct_name
        
        # 提取结构体名前缀（去除动态前缀）
        type_suffix = return_type.replace(struct_prefix, '', 1)
        
        methods.append({
            'method_name': method_name,
            'return_type': return_type,
            'event_name': event_name,
            'solidity_signature': solidity_signature,
            'struct_prefix': type_suffix,
            'receiver_var': receiver_var,
            'struct_name': struct_name
        })
    
    return methods, struct_prefix


def extract_pack_methods(go_file_path, struct_name, receiver_var):
    """提取所有PackXXX方法信息"""
    try:
        with open(go_file_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except FileNotFoundError:
        print(f"错误：找不到文件 {go_file_path}")
        sys.exit(1)
    
    # 匹配PackXXX方法的模式，同时提取Solidity函数签名和methodID（通用化结构体前缀）
    pattern = rf'// Pack(\w+) is the Go binding.*?method with ID (0x[a-fA-F0-9]+).*?// Solidity: function ([^\n]+)\nfunc \({re.escape(receiver_var)} \*{re.escape(struct_name)}\) Pack\1\([^)]*\) \[\]byte \{{\s*enc, err := {re.escape(receiver_var)}\.abi\.Pack\("([^"]+)"\)'
    
    methods = []
    for match in re.finditer(pattern, content, re.DOTALL):
        pack_method_name = match.group(1)  # 例如: AuthorizedSigner
        method_id = match.group(2)  # 例如: 0xc771909c
        solidity_signature = match.group(3).strip()  # 例如: authorizedSigner() view returns(address)
        abi_method_name = match.group(4)  # 例如: authorizedSigner
        
        methods.append({
            'pack_method_name': pack_method_name,
            'method_id': method_id,
            'solidity_signature': solidity_signature,
            'abi_method_name': abi_method_name
        })
    
    return methods


def clean_solidity_signature(signature):
    """清理Solidity事件签名，去除indexed关键字"""
    # 去除 indexed 关键字
    cleaned = signature.replace(' indexed', '')
    # 去除多余的空格
    cleaned = ' '.join(cleaned.split())
    return cleaned


def calculate_topic0_hash(event_signature):
    """使用cast sig-event计算事件的topic0 hash"""
    try:
        result = subprocess.run(
            ['cast', 'sig-event', event_signature],
            capture_output=True,
            text=True,
            check=True
        )
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        print(f"错误：无法计算事件签名 hash: {event_signature}")
        print(f"cast 命令输出: {e.stderr}")
        return None
    except FileNotFoundError:
        print("错误：cast 命令未找到，请确保已安装 Foundry")
        return None


def generate_unpack_event(methods, struct_name, receiver_var):
    """生成UnpackEvent方法的Go代码"""
    code_lines = [
        "// UnpackEvent unpacks event log based on topic0.",
        f"func ({receiver_var} *{struct_name}) UnpackEvent(log *types.Log) (interface {{",
        f"\tContractEventName() string",
        f"\tTopic0() common.Hash",
        f"}}, error) {{",
        "\tvar mismatch = errors.New(\"event signature mismatch\")",
        "\tif len(log.Topics) == 0 {",
        "\t\treturn nil, mismatch",
        "\t}",
        "\ttopic0 := log.Topics[0]",
    ]
    
    for method in methods:
        code_lines.extend([
            f"\tif topic0 == {receiver_var}.abi.Events[\"{method['event_name']}\"].ID {{",
            f"\t\treturn {receiver_var}.{method['method_name']}(log)",
            f"\t}}",
        ])
    
    code_lines.extend([
        "\treturn nil, mismatch",
        "}"
    ])
    
    return "\n".join(code_lines)


def generate_method_id_methods(pack_methods, struct_name, receiver_var):
    """生成函数MethodID方法"""
    code_lines = []
    
    for method in pack_methods:
        method_name = f"{method['pack_method_name']}MethodID"
        
        # 提取methodID的后8位（去掉0x前缀）
        method_id_hex = method['method_id'][2:]  # 去掉0x
        
        code_lines.extend([
            f"// {method_name} returns the method ID for {method['abi_method_name']} ({method['method_id']}).",
            f"//",
            f"// Solidity: function {method['solidity_signature']}",
            f"func ({receiver_var} *{struct_name}) {method_name}() [4]byte {{",
            f"\treturn [4]byte{{0x{method_id_hex[:2]}, 0x{method_id_hex[2:4]}, 0x{method_id_hex[4:6]}, 0x{method_id_hex[6:8]}}}",
            f"}}",
            ""
        ])
    
    return "\n".join(code_lines[:-1])  # 去除最后一个空行


def generate_event_topic0_methods(methods, struct_name):
    """生成事件Topic0方法"""
    code_lines = []
    
    for method in methods:
        # 清理Solidity签名
        cleaned_signature = clean_solidity_signature(method['solidity_signature'])
        
        # 计算topic0 hash
        topic0_hash = calculate_topic0_hash(cleaned_signature)
        if topic0_hash is None:
            continue
            
        method_name = f"{struct_name}{method['struct_prefix']}Topic0"
        
        code_lines.extend([
            f"// {method_name} returns the hash of the event signature.",
            f"//",
            f"// Solidity: event {cleaned_signature}",
            f"func {method_name}() common.Hash {{",
            f"\treturn common.HexToHash(\"{topic0_hash}\")",
            f"}}",
            ""
        ])
    
    return "\n".join(code_lines[:-1])  # 去除最后一个空行


def generate_struct_topic0_methods(methods, struct_name):
    """生成结构体的Topic0方法"""
    code_lines = []
    
    for method in methods:
        # 清理Solidity签名
        cleaned_signature = clean_solidity_signature(method['solidity_signature'])
        
        # 计算topic0 hash
        topic0_hash = calculate_topic0_hash(cleaned_signature)
        if topic0_hash is None:
            continue
            
        struct_type = method['return_type']  # 例如: T2AuthorizedSignerUpdated
        
        code_lines.extend([
            f"// Topic0 returns the hash of the event signature.",
            f"//",
            f"// Solidity: event {cleaned_signature}",
            f"func ({struct_type}) Topic0() common.Hash {{",
            f"\treturn common.HexToHash(\"{topic0_hash}\")",
            f"}}",
            ""
        ])
    
    return "\n".join(code_lines[:-1])  # 去除最后一个空行


def check_existing_methods(content, struct_name, receiver_var):
    """检查已存在的方法，返回位置信息"""
    # 检查UnpackEvent方法
    unpack_pattern = rf'// UnpackEvent.*?^func \({re.escape(receiver_var)} \*{re.escape(struct_name)}\) UnpackEvent.*?^}}'
    unpack_match = re.search(unpack_pattern, content, re.DOTALL | re.MULTILINE)
    
    # 检查Topic0方法区域（从第一个Topic0方法开始到最后一个结束）
    topic0_pattern = rf'// {re.escape(struct_name)}\w+Topic0 returns.*?^func {re.escape(struct_name)}\w+Topic0\(\).*?^}}'
    topic0_matches = list(re.finditer(topic0_pattern, content, re.DOTALL | re.MULTILINE))
    
    # 检查MethodID方法区域
    methodid_pattern = rf'// \w+MethodID returns.*?^func \({re.escape(receiver_var)} \*{re.escape(struct_name)}\) \w+MethodID\(\).*?^}}'
    methodid_matches = list(re.finditer(methodid_pattern, content, re.DOTALL | re.MULTILINE))
    
    # 检查结构体Topic0方法区域
    struct_topic0_pattern = rf'// Topic0 returns the hash of the event signature\.\s*//\s*// Solidity: event.*?^func \({re.escape(struct_name)}\w+\) Topic0\(\).*?^}}'
    struct_topic0_matches = list(re.finditer(struct_topic0_pattern, content, re.DOTALL | re.MULTILINE))
    
    unpack_pos = None
    topic0_pos = None
    methodid_pos = None
    struct_topic0_pos = None
    
    if unpack_match:
        unpack_pos = (unpack_match.start(), unpack_match.end())
    
    if topic0_matches:
        # 找到所有Topic0方法的范围
        first_start = topic0_matches[0].start()
        last_end = topic0_matches[-1].end()
        topic0_pos = (first_start, last_end)
    
    if methodid_matches:
        # 找到所有MethodID方法的范围
        first_start = methodid_matches[0].start()
        last_end = methodid_matches[-1].end()
        methodid_pos = (first_start, last_end)
    
    if struct_topic0_matches:
        # 找到所有结构体Topic0方法的范围
        first_start = struct_topic0_matches[0].start()
        last_end = struct_topic0_matches[-1].end()
        struct_topic0_pos = (first_start, last_end)
    
    return unpack_pos, topic0_pos, methodid_pos, struct_topic0_pos


def write_methods_to_file(go_file_path, unpack_code, topic0_code, methodid_code, struct_topic0_code, struct_name, receiver_var):
    """将生成的方法写入到Go文件中，如果已存在则覆盖"""
    try:
        with open(go_file_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except FileNotFoundError:
        print(f"错误：找不到文件 {go_file_path}")
        return False
    
    # 检查已存在的方法
    unpack_pos, topic0_pos, methodid_pos, struct_topic0_pos = check_existing_methods(content, struct_name, receiver_var)
    
    new_content = content
    
    # 处理UnpackEvent方法
    if unpack_pos is not None:
        new_content = new_content[:unpack_pos[0]] + unpack_code + new_content[unpack_pos[1]:]
#         print("检测到已存在的UnpackEvent方法，正在覆盖...")
    else:
        if not new_content.endswith('\n'):
            new_content = new_content + '\n\n' + unpack_code
        else:
            new_content = new_content + '\n' + unpack_code
#         print("在文件末尾添加UnpackEvent方法...")
    
    # 重新检查位置（因为内容可能已经变化）
    unpack_pos_new, topic0_pos_new, methodid_pos_new, struct_topic0_pos_new = check_existing_methods(new_content, struct_name, receiver_var)
    
    # 处理Topic0方法
    if topic0_pos_new is not None:
        new_content = new_content[:topic0_pos_new[0]] + topic0_code + new_content[topic0_pos_new[1]:]
#         print("检测到已存在的Topic0方法，正在覆盖...")
    else:
        if not new_content.endswith('\n'):
            new_content = new_content + '\n\n' + topic0_code
        else:
            new_content = new_content + '\n' + topic0_code
#         print("在文件末尾添加Topic0方法...")
    
    # 重新检查位置
    unpack_pos_final, topic0_pos_final, methodid_pos_final, struct_topic0_pos_final = check_existing_methods(new_content, struct_name, receiver_var)
    
    # 处理MethodID方法
    if methodid_code and methodid_pos_final is not None:
        new_content = new_content[:methodid_pos_final[0]] + methodid_code + new_content[methodid_pos_final[1]:]
#         print("检测到已存在的MethodID方法，正在覆盖...")
    elif methodid_code:
        if not new_content.endswith('\n'):
            new_content = new_content + '\n\n' + methodid_code
        else:
            new_content = new_content + '\n' + methodid_code
#         print("在文件末尾添加MethodID方法...")
    
    # 最后检查位置
    unpack_pos_last, topic0_pos_last, methodid_pos_last, struct_topic0_pos_last = check_existing_methods(new_content, struct_name, receiver_var)
    
    # 处理结构体Topic0方法
    if struct_topic0_code and struct_topic0_pos_last is not None:
        new_content = new_content[:struct_topic0_pos_last[0]] + struct_topic0_code + new_content[struct_topic0_pos_last[1]:]
#         print("检测到已存在的结构体Topic0方法，正在覆盖...")
    elif struct_topic0_code:
        if not new_content.endswith('\n'):
            new_content = new_content + '\n\n' + struct_topic0_code
        else:
            new_content = new_content + '\n' + struct_topic0_code
#         print("在文件末尾添加结构体Topic0方法...")
    
    try:
        with open(go_file_path, 'w', encoding='utf-8') as f:
            f.write(new_content)
        return True
    except Exception as e:
        print(f"写入文件时出错: {e}")
        return False


def main():
    if len(sys.argv) != 2:
        print("用法: python3 generate_unpack_event.py <go_file_path>")
        sys.exit(1)
    
    go_file_path = sys.argv[1]
    
#     print(f"正在分析文件: {go_file_path}")
    
    # 提取UnpackXXXEvent方法
    unpack_methods, struct_prefix = extract_unpack_methods(go_file_path)
    
    if not unpack_methods:
        print("未找到任何UnpackXXXEvent方法")
        sys.exit(1)
    
    # 获取结构体信息
    struct_name = unpack_methods[0]['struct_name']
    receiver_var = unpack_methods[0]['receiver_var']
    
#     print(f"检测到结构体: {struct_name}, 接收器变量: {receiver_var}")
#     print(f"找到 {len(unpack_methods)} 个UnpackXXXEvent方法:")
#     for method in unpack_methods:
#         print(f"  - {method['method_name']} -> {method['return_type']} (event: {method['event_name']})")
#         print(f"    Solidity: {method['solidity_signature']}")
    
    # 提取PackXXX方法
    pack_methods = extract_pack_methods(go_file_path, struct_name, receiver_var)
    
#     if pack_methods:
#         print(f"\n找到 {len(pack_methods)} 个PackXXX方法:")
#         for method in pack_methods:
#             print(f"  - Pack{method['pack_method_name']} -> {method['method_id']}")
#             print(f"    Solidity: function {method['solidity_signature']}")
#     else:
#         print("\n未找到任何PackXXX方法")
    
    # 生成UnpackEvent代码
    unpack_code = generate_unpack_event(unpack_methods, struct_name, receiver_var)
    
    # 生成Topic0方法代码
#     print("\n正在生成Topic0方法...")
    topic0_code = generate_event_topic0_methods(unpack_methods, struct_name)
    
    # 生成MethodID方法代码
    methodid_code = ""
    if pack_methods:
#         print("正在生成MethodID方法...")
        methodid_code = generate_method_id_methods(pack_methods, struct_name, receiver_var)
    
    # 生成结构体Topic0方法代码
#     print("正在生成结构体Topic0方法...")
    struct_topic0_code = generate_struct_topic0_methods(unpack_methods, struct_name)
    
    # 写入到原文件
    if write_methods_to_file(go_file_path, unpack_code, topic0_code, methodid_code, struct_topic0_code, struct_name, receiver_var):
        print(f"\n✓ 所有方法已成功写入到 {go_file_path}")
    else:
        print("\n✗ 写入文件失败")
        sys.exit(1)
    
#     print("\n生成的UnpackEvent代码:")
#     print("=" * 50)
#     print(unpack_code)
#     print("=" * 50)
#
#     print("\n生成的Topic0方法代码:")
#     print("=" * 50)
#     print(topic0_code)
#     print("=" * 50)
    
#     if methodid_code:
#         print("\n生成的MethodID方法代码:")
#         print("=" * 50)
#         print(methodid_code)
#         print("=" * 50)
#
#     print("\n生成的结构体Topic0方法代码:")
#     print("=" * 50)
#     print(struct_topic0_code)
#     print("=" * 50)


if __name__ == "__main__":
    main()
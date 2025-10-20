#!/usr/bin/env python3
"""
Go事件解析器生成工具

使用方法:
python3 ./scripts/generate_unpack_event.py <go_file_path>
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


def parse_go_params(param_str):
    """解析Go函数参数字符串，返回参数列表"""
    if not param_str.strip():
        return []
    
    params = []
    # 处理参数列表，支持多个参数
    # 例如: "to common.Address, value *big.Int"
    parts = param_str.split(',')
    for part in parts:
        part = part.strip()
        if not part:
            continue
        # 分割参数名和类型
        tokens = part.split()
        if len(tokens) >= 2:
            param_name = tokens[0]
            param_type = ' '.join(tokens[1:])
            params.append({'name': param_name, 'type': param_type})
    
    return params


def parse_solidity_params(signature):
    """从Solidity函数签名中提取参数类型
    例如: "transfer(address to, uint256 value) returns(bool)" -> [('address', 'to'), ('uint256', 'value')]
    """
    # 提取函数参数部分
    match = re.search(r'\(([^)]*)\)', signature)
    if not match:
        return []
    
    params_str = match.group(1).strip()
    if not params_str:
        return []
    
    params = []
    # 分割参数
    for param in params_str.split(','):
        param = param.strip()
        if not param:
            continue
        # 分割类型和名称
        parts = param.split()
        if len(parts) >= 2:
            param_type = parts[0]
            param_name = parts[1]
            params.append((param_type, param_name))
        elif len(parts) == 1:
            # 只有类型，没有名称
            params.append((parts[0], ''))
    
    return params


def extract_pack_methods(go_file_path, struct_name, receiver_var):
    """提取所有PackXXX方法信息（包含详细参数信息）"""
    try:
        with open(go_file_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except FileNotFoundError:
        print(f"错误：找不到文件 {go_file_path}")
        sys.exit(1)
    
    # 匹配PackXXX方法的模式，同时提取Solidity函数签名、methodID和Go参数
    pattern = rf'// Pack(\w+) is the Go binding.*?method with ID (0x[a-fA-F0-9]+).*?// Solidity: function ([^\n]+)\nfunc \({re.escape(receiver_var)} \*{re.escape(struct_name)}\) Pack\1\(([^)]*)\) \[\]byte \{{\s*enc, err := {re.escape(receiver_var)}\.abi\.Pack\("([^"]+)"'
    
    methods = []
    for match in re.finditer(pattern, content, re.DOTALL):
        pack_method_name = match.group(1)  # 例如: Transfer
        method_id = match.group(2)  # 例如: 0xa9059cbb
        solidity_signature = match.group(3).strip()  # 例如: transfer(address to, uint256 value) returns(bool)
        go_params_str = match.group(4).strip()  # 例如: to common.Address, value *big.Int
        abi_method_name = match.group(5)  # 例如: transfer
        
        # 解析Go参数
        go_params = parse_go_params(go_params_str)
        
        # 解析Solidity参数
        solidity_params = parse_solidity_params(solidity_signature)
        
        methods.append({
            'pack_method_name': pack_method_name,
            'method_id': method_id,
            'solidity_signature': solidity_signature,
            'abi_method_name': abi_method_name,
            'go_params': go_params,
            'solidity_params': solidity_params
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


def get_go_zero_value(go_type):
    """获取Go类型的零值"""
    if go_type == "common.Address":
        return "common.Address{}"
    elif go_type == "*big.Int":
        return "nil"
    elif go_type == "bool":
        return "false"
    elif go_type == "uint8":
        return "0"
    elif go_type == "string":
        return "\"\""
    else:
        return "nil"


def is_pointer_type(go_type):
    """判断Go类型是否应该是指针类型"""
    # 已经是指针的类型
    if go_type.startswith("*"):
        return True
    # 基本类型和特殊类型
    if go_type in ["common.Address", "bool", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "string", "[]byte"]:
        return False
    # 数组/切片类型
    if go_type.startswith("[]") or go_type.startswith("["):
        return False
    # 其他类型默认应该是指针（如结构体）
    return True


def generate_unpack_input_method(method, struct_name, receiver_var):
    """生成单个UnpackInputXXX方法"""
    method_name = f"UnpackInput{method['pack_method_name']}"
    abi_method_name = method['abi_method_name']
    go_params = method['go_params']
    
    code_lines = [
        f"// {method_name} unpacks the input data for the {abi_method_name} method.",
        f"//",
        f"// Solidity: function {method['solidity_signature']}",
    ]
    
    # 构建函数签名
    if go_params:
        # 有参数
        # 检查每个参数类型，如果不是指针类型但应该是指针的，则添加指针
        adjusted_params = []
        for p in go_params:
            param_type = p['type']
            # 如果类型不是指针但应该是指针（如结构体），则添加指针
            if is_pointer_type(param_type) and not param_type.startswith("*"):
                adjusted_params.append({'name': p['name'], 'type': f"*{param_type}"})
            else:
                adjusted_params.append(p)
        
        return_params = ", ".join([f"{p['name']} {p['type']}" for p in adjusted_params])
        zero_values = ", ".join([get_go_zero_value(p['type']) for p in adjusted_params])
        
        code_lines.append(f"func ({receiver_var} *{struct_name}) {method_name}(callData []byte) ({return_params}, err error) {{")
        code_lines.append(f"\tmethod, ok := {receiver_var}.abi.Methods[\"{abi_method_name}\"]")
        code_lines.append(f"\tif !ok {{")
        code_lines.append(f"\t\treturn {zero_values}, errors.New(\"method '{abi_method_name}' not found\")")
        code_lines.append(f"\t}}")
        code_lines.append(f"\tif len(callData) < 4 || !bytes.Equal(callData[:4], method.ID[:4]) {{")
        code_lines.append(f"\t\treturn {zero_values}, errors.New(\"method signature mismatch\")")
        code_lines.append(f"\t}}")
        code_lines.append(f"\targuments, err := method.Inputs.Unpack(callData[4:])")
        code_lines.append(f"\tif err != nil {{")
        code_lines.append(f"\t\treturn {zero_values}, err")
        code_lines.append(f"\t}}")
        
        # 生成参数解析代码
        for i, param in enumerate(adjusted_params):
            param_name = param['name']
            param_type = param['type']
            
            if param_type == "common.Address":
                code_lines.append(f"\t{param_name} = *abi.ConvertType(arguments[{i}], new(common.Address)).(*common.Address)")
            elif param_type == "*big.Int":
                code_lines.append(f"\t{param_name} = abi.ConvertType(arguments[{i}], new(big.Int)).(*big.Int)")
            elif param_type == "bool":
                code_lines.append(f"\t{param_name} = *abi.ConvertType(arguments[{i}], new(bool)).(*bool)")
            elif param_type == "uint8":
                code_lines.append(f"\t{param_name} = *abi.ConvertType(arguments[{i}], new(uint8)).(*uint8)")
            elif param_type == "string":
                code_lines.append(f"\t{param_name} = *abi.ConvertType(arguments[{i}], new(string)).(*string)")
            elif param_type == "[]byte":
                code_lines.append(f"\t{param_name} = abi.ConvertType(arguments[{i}], new([]byte)).([]byte)")
            elif param_type.startswith("*"):
                # 指针类型（包括结构体指针）
                inner_type = param_type[1:]  # 去掉*
                code_lines.append(f"\t{param_name} = abi.ConvertType(arguments[{i}], new({inner_type})).(*{inner_type})")
            else:
                # 默认处理
                code_lines.append(f"\t{param_name} = abi.ConvertType(arguments[{i}], new({param_type})).({param_type})")
        
        code_lines.append(f"\treturn {', '.join([p['name'] for p in adjusted_params])}, nil")
    else:
        # 无参数
        code_lines.append(f"func ({receiver_var} *{struct_name}) {method_name}(callData []byte) error {{")
        code_lines.append(f"\tmethod, ok := {receiver_var}.abi.Methods[\"{abi_method_name}\"]")
        code_lines.append(f"\tif !ok {{")
        code_lines.append(f"\t\treturn errors.New(\"method '{abi_method_name}' not found\")")
        code_lines.append(f"\t}}")
        code_lines.append(f"\tif len(callData) < 4 || !bytes.Equal(callData[:4], method.ID[:4]) {{")
        code_lines.append(f"\t\treturn errors.New(\"method signature mismatch\")")
        code_lines.append(f"\t}}")
        code_lines.append(f"\treturn nil")
    
    code_lines.append("}")
    
    return "\n".join(code_lines)


def generate_unpack_input_methods(pack_methods, struct_name, receiver_var):
    """生成所有UnpackInputXXX方法"""
    code_lines = []
    
    for method in pack_methods:
        code_lines.append(generate_unpack_input_method(method, struct_name, receiver_var))
        code_lines.append("")
    
    return "\n".join(code_lines[:-1])  # 去除最后一个空行


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
    # 检查UnpackEvent方法（需要匹配多行返回类型）
    unpack_pattern = rf'// UnpackEvent unpacks event log based on topic0\.\s*^func \({re.escape(receiver_var)} \*{re.escape(struct_name)}\) UnpackEvent\(log \*types\.Log\).*?^\}}\s*$'
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
    
    # 检查UnpackInputXXX方法区域
    unpack_input_pattern = rf'// UnpackInput\w+ unpacks the input data.*?^func \({re.escape(receiver_var)} \*{re.escape(struct_name)}\) UnpackInput\w+\(callData \[\]byte\).*?^}}'
    unpack_input_matches = list(re.finditer(unpack_input_pattern, content, re.DOTALL | re.MULTILINE))
    
    unpack_pos = None
    topic0_pos = None
    methodid_pos = None
    struct_topic0_pos = None
    unpack_input_pos = None
    
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
    
    if unpack_input_matches:
        # 找到所有UnpackInputXXX方法的范围
        first_start = unpack_input_matches[0].start()
        last_end = unpack_input_matches[-1].end()
        unpack_input_pos = (first_start, last_end)
    
    return unpack_pos, topic0_pos, methodid_pos, struct_topic0_pos, unpack_input_pos


def write_methods_to_file(go_file_path, unpack_code, topic0_code, methodid_code, struct_topic0_code, unpack_input_code, struct_name, receiver_var):
    """将生成的方法写入到Go文件中，如果已存在则覆盖"""
    try:
        with open(go_file_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except FileNotFoundError:
        print(f"错误：找不到文件 {go_file_path}")
        return False
    
    # 检查已存在的方法
    unpack_pos, topic0_pos, methodid_pos, struct_topic0_pos, unpack_input_pos = check_existing_methods(content, struct_name, receiver_var)
    
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
    unpack_pos_new, topic0_pos_new, methodid_pos_new, struct_topic0_pos_new, unpack_input_pos_new = check_existing_methods(new_content, struct_name, receiver_var)
    
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
    unpack_pos_final, topic0_pos_final, methodid_pos_final, struct_topic0_pos_final, unpack_input_pos_final = check_existing_methods(new_content, struct_name, receiver_var)
    
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
    
    # 再次检查位置
    unpack_pos_temp, topic0_pos_temp, methodid_pos_temp, struct_topic0_pos_temp, unpack_input_pos_temp = check_existing_methods(new_content, struct_name, receiver_var)
    
    # 处理UnpackInputXXX方法
    if unpack_input_code and unpack_input_pos_temp is not None:
        new_content = new_content[:unpack_input_pos_temp[0]] + unpack_input_code + new_content[unpack_input_pos_temp[1]:]
#         print("检测到已存在的UnpackInput方法，正在覆盖...")
    elif unpack_input_code:
        if not new_content.endswith('\n'):
            new_content = new_content + '\n\n' + unpack_input_code
        else:
            new_content = new_content + '\n' + unpack_input_code
#         print("在文件末尾添加UnpackInput方法...")
    
    # 最后检查位置
    unpack_pos_last, topic0_pos_last, methodid_pos_last, struct_topic0_pos_last, unpack_input_pos_last = check_existing_methods(new_content, struct_name, receiver_var)
    
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
#             print(f"    Go参数: {method['go_params']}")
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
    
    # 生成UnpackInputXXX方法代码
    unpack_input_code = ""
    if pack_methods:
#         print("正在生成UnpackInput方法...")
        unpack_input_code = generate_unpack_input_methods(pack_methods, struct_name, receiver_var)
    
    # 生成结构体Topic0方法代码
#     print("正在生成结构体Topic0方法...")
    struct_topic0_code = generate_struct_topic0_methods(unpack_methods, struct_name)
    
    # 写入到原文件
    if write_methods_to_file(go_file_path, unpack_code, topic0_code, methodid_code, struct_topic0_code, unpack_input_code, struct_name, receiver_var):
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
#     if unpack_input_code:
#         print("\n生成的UnpackInput方法代码:")
#         print("=" * 50)
#         print(unpack_input_code)
#         print("=" * 50)
#
#     print("\n生成的结构体Topic0方法代码:")
#     print("=" * 50)
#     print(struct_topic0_code)
#     print("=" * 50)


if __name__ == "__main__":
    main()
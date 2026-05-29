
build:
	go run ../abigen -v2 -pkg aaa --abi Escrow.json --type "EscrowPkg" --alias eip712Domain="EscrowEip712Domain" --out "EscrowPack.go"

build2:
	go run ../abigen -pkg aaa --abi "./Escrow.json" --type "Escrow" --alias eip712Domain="EscrowEip712Domain" --out "./Escrow.go"
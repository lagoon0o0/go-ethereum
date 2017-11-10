package main

import (
"fmt"
//"encoding/hex"
"github.com/ethereum/go-ethereum/rlp"
//"bytes"
"encoding/hex"
"github.com/ethereum/go-ethereum/ethdb"
//"github.com/ethereum/go-ethereum/common"
//"github.com/ethereum/go-ethereum/core/state"
"github.com/ethereum/go-ethereum/consensus/ethash"
"github.com/ethereum/go-ethereum/core"
"github.com/ethereum/go-ethereum/core/vm"
"github.com/ethereum/go-ethereum/common"
"bytes"
"github.com/ethereum/go-ethereum/crypto"
"reflect"
"github.com/ethereum/go-ethereum/core/types"
"github.com/ethereum/go-ethereum/trie"
"github.com/ethereum/go-ethereum/node"
"github.com/ethereum/go-ethereum/eth"
"github.com/ethereum/go-ethereum/log"
"time"

"github.com/ethereum/go-ethereum/les"

)

type Result struct{
	Proof []rlp.RawValue
	Key	  []byte

}

func Pretty_print(str string, arr interface{} ) {
	if reflect.TypeOf(arr).Kind() == reflect.Slice {
		t := reflect.ValueOf(arr)
		fmt.Println(str, t.Len(), arr)
		for i := 0; i < t.Len(); i++ {
			fmt.Printf("%x ",t.Index(i))
		}
		fmt.Println("")
	}

}

func print_proof(proof []rlp.RawValue)  {
	enc,_ := rlp.EncodeToBytes(proof)
	fmt.Println("The proof: ")
	for _, node := range proof {

		Pretty_print("encoded node: ",node)

		r := bytes.NewReader(node)
		decoded := new([][]byte)
		_ = rlp.Decode(r, &decoded)
		Pretty_print("decoded node: ", *decoded)
	}
	Pretty_print("\nProof encoding: ", enc)

}
func request_account(chaindb_dir,state_root, account string) (string, []string) {
	if state_root[:2] == "0x" {
		state_root = state_root[2:]
	}

	// load the block chain
	chainDb, err := ethdb.NewLDBDatabase(chaindb_dir,1024, 64)
	engine := ethash.NewFaker()
	config, _, err := core.SetupGenesisBlock(chainDb, nil)
	vmcfg := vm.Config{EnablePreimageRecording: false}
	chain, err := core.NewBlockChain(chainDb, config, engine, vmcfg)

	// get the root state
	decoded, err := hex.DecodeString(state_root)
	if err != nil {
		log.Info("Error")
	}
	h := common.BytesToHash(decoded)
	sta, err := chain.StateAt(h)
	t := sta.GetTrie()

	// get the proof
	addr := common.HexToAddress(account)
	proof := t.GetProof(addr[:])
	enc,_ := rlp.EncodeToBytes(proof)


	veri_enc,_ := t.VerifyProof(addr[:], proof)

	var value [][]byte
	if err := rlp.DecodeBytes(veri_enc, &value); err != nil {
		fmt.Println("Error: ", err)
	}
	//Pretty_print("decoded node: ", value)

	//fmt.Println("key:",crypto.Keccak256Hash(addr.Bytes()).Hex())

	tuple := []string{
		hex.EncodeToString(value[0]),
		hex.EncodeToString(value[1]),
		hex.EncodeToString(value[2]),
		hex.EncodeToString(value[3])}

	ret := []interface{}{proof, addr.Bytes()}
	enc, _ = rlp.EncodeToBytes(ret)
	chainDb.Close()
	return hex.EncodeToString(enc), tuple
}

func request_storage(chaindb_dir,state_root, account string, key []byte)(string){
	if state_root[:2] == "0x" {
		state_root = state_root[2:]
	}
	// load the block chain
	chainDb, err := ethdb.NewLDBDatabase(chaindb_dir,1024, 64)
	engine := ethash.NewFaker()
	config, _, err := core.SetupGenesisBlock(chainDb, nil)
	vmcfg := vm.Config{EnablePreimageRecording: false}
	chain, err := core.NewBlockChain(chainDb, config, engine, vmcfg)

	// get the root state
	decoded, err := hex.DecodeString(state_root)
	if err != nil {
		log.Info("Error")
	}
	h := common.BytesToHash(decoded)
	sta, err := chain.StateAt(h)


	addr := common.HexToAddress(account)

	t := sta.GetStorageSecureTrie(addr)

	fmt.Println(key,common.Bytes2Hex(key))

	proof := t.GetProof(key)

	enc,_ := rlp.EncodeToBytes(proof)
	fmt.Println("The proof: ")
	for _, node := range proof {
		Pretty_print("encoded node: ",node)
		fmt.Println("Hash:",crypto.Keccak256Hash(node).Hex())
		r := bytes.NewReader(node)
		decoded := new([][]byte)
		_ = rlp.Decode(r, &decoded)
		Pretty_print("decoded node: ", *decoded)
	}
	Pretty_print("\nProof encoding: ", enc)

	enc,_ = t.VerifyProof(key, proof)
	if err != nil {
		fmt.Println("Error:",err)
	}
	Pretty_print("\nVerify:", enc)


	var value []byte
	if err := rlp.DecodeBytes(enc, &value); err != nil {
		fmt.Println("Error: ", err)
	}
	Pretty_print("decoded node: ", value)



	ret := []interface{}{proof, key}
	enc, _ = rlp.EncodeToBytes(ret)
	Pretty_print("ret: ", ret)
	fmt.Println("")
	fmt.Println("enc: ",enc)
	fmt.Println("enc: ",hex.EncodeToString(enc))
	chainDb.Close()
	return hex.EncodeToString(enc)
}

func request_transaction(chaindb_dir string, tx_hash string, block_number uint64)(int, string) {
	// load the block chain
	chainDb, _ := ethdb.NewLDBDatabase(chaindb_dir,1024, 64)
	engine := ethash.NewFaker()
	config, _, _ := core.SetupGenesisBlock(chainDb, nil)
	vmcfg := vm.Config{EnablePreimageRecording: false}
	chain, _ := core.NewBlockChain(chainDb, config, engine, vmcfg)

	block := chain.GetBlockByNumber(block_number)

	receipt, blockHash, blockNumber, receiptIndex := core.GetReceipt(chainDb,common.HexToHash(tx_hash))

	fmt.Println(receipt)
	fmt.Println(blockHash)
	fmt.Println(blockNumber)
	fmt.Println(receiptIndex)




	tx_list := block.Transactions()
	tx_idx := -1
	var tx *types.Transaction
	for idx , transaction := range block.Transactions() {

		if transaction.Hash() == common.HexToHash(tx_hash) {
			tx_idx = idx
			tx = transaction
			break
		}
	}
	if tx_idx == -1 {
		fmt.Println("transaction not found")
		return -1, ""
	}
	fmt.Println(tx)

	keybuf := new(bytes.Buffer)
	t := new(trie.Trie)
	for i := 0; i < tx_list.Len(); i++ {
		keybuf.Reset()
		rlp.Encode(keybuf, uint(i))
		t.Update(keybuf.Bytes(), tx_list.GetRlp(i))
	}
	tx_root := t.Hash()
	fmt.Println(tx_root.Hex())

	rlp.Encode(keybuf, uint(tx_idx))

	proof := t.Prove(keybuf.Bytes())

	enc,_ := rlp.EncodeToBytes(proof)
	/*
	fmt.Println("The proof: ")
	for _, node := range proof {
		Pretty_print("encoded node: ",node)
		fmt.Println("Hash:",crypto.Keccak256Hash(node).Hex())
		r := bytes.NewReader(node)
		decoded := new([][]byte)
		_ = rlp.Decode(r, &decoded)
		Pretty_print("decoded node: ", *decoded)
	}
	Pretty_print("\nProof encoding: ", enc)
	*/
	chainDb.Close()
	return tx_idx, hex.EncodeToString(enc)
}

func request_proof(chaindb_dir, state_root,  account string, user string, token string, pos string) (string, string){
	if state_root[:2] == "0x" {
		state_root = state_root[2:]
	}
	if user[:2] == "0x" {
		user = user[2:]
	}
	// load the block chain
	chainDb, err := ethdb.NewLDBDatabase(chaindb_dir,1024, 64)
	engine := ethash.NewFaker()
	config, _, err := core.SetupGenesisBlock(chainDb, nil)
	vmcfg := vm.Config{EnablePreimageRecording: false}
	chain, err := core.NewBlockChain(chainDb, config, engine, vmcfg)



	// get the root state
	decoded, err := hex.DecodeString(state_root)
	if err != nil {
		log.Info("fuck")
	}
	h := common.BytesToHash(decoded)
	sta, err := chain.StateAt(h)
	state_trie := sta.GetTrie()

	// get the proof
	addr := common.HexToAddress(account)
	account_proof := state_trie.GetProof(addr[:])


	storage_trie := sta.GetStorageSecureTrie(addr)
	key := crypto.Keccak256Hash(
		common.Hex2BytesFixed(user, 32),
		crypto.Keccak256Hash(
			common.Hex2BytesFixed(token, 32),
			common.Hex2BytesFixed(pos, 32)).Bytes()).Bytes()

	fmt.Println(crypto.Keccak256Hash(
		common.Hex2BytesFixed(token, 32),
		common.Hex2BytesFixed(pos, 32)).Hex())
	fmt.Println(common.Bytes2Hex(common.Hex2BytesFixed(user, 32)))
	fmt.Println(crypto.Keccak256Hash(
		common.Hex2BytesFixed(user, 32),
		crypto.Keccak256Hash(
			common.Hex2BytesFixed(token, 32),
			common.Hex2BytesFixed(pos, 32)).Bytes()).Hex())
	storage_proof := storage_trie.GetProof(key)

	value, err := storage_trie.VerifyProof(key,storage_proof)
	if err != nil {
		fmt.Println("Verification Error:",err)
	}
	//fmt.Println("Value: ",value)
	fmt.Println("user",hex.EncodeToString(common.Hex2BytesFixed(user, 32)))
	ret := []interface{}{
		account_proof,
		addr.Bytes(),
		storage_proof,
		common.Hex2BytesFixed(user, 32),
		common.Hex2BytesFixed(token, 32),
		common.Hex2BytesFixed(pos, 32)}
	enc, _ := rlp.EncodeToBytes(ret)
	fmt.Println(hex.EncodeToString(enc))
	chainDb.Close()
	return hex.EncodeToString(enc), hex.EncodeToString(value)
}

func request_header(chaindb_dir string, number uint64)(string, *types.Header){
	// load the block chain
	chainDb, _ := ethdb.NewLDBDatabase(chaindb_dir,1024, 64)
	engine := ethash.NewFaker()
	config, _, _ := core.SetupGenesisBlock(chainDb, nil)
	vmcfg := vm.Config{EnablePreimageRecording: false}
	chain, _ := core.NewBlockChain(chainDb, config, engine, vmcfg)


	header := chain.GetHeaderByNumber(number)

	//fmt.Println(header)
	fmt.Println("Miner Hash:",header.HashNoNonce().Hex())

	enc, _ := rlp.EncodeToBytes(header)
	//fmt.Println(hex.EncodeToString(enc))
	chainDb.Close()
	return hex.EncodeToString(enc), header

}

func get_database() {
	config := node.DefaultConfig
	config.HTTPModules = []string{"db","eth","net","web3","personal"}
	config.Name = "geth"
	config.DataDir = "/Users/lagoon/Library/Ethereum"
	config.HTTPHost = "localhost"
	config.HTTPPort = 1234
	stack, err := node.New(&config)

	if err != nil {
		fmt.Println("error: ", err)
	}

	cfg := eth.DefaultConfig
	err = stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		fullNode, err := eth.New(ctx, &cfg)

		if fullNode != nil && cfg.LightServ > 0 {
			ls, _ := les.NewLesServer(fullNode, &cfg)
			fullNode.AddLesServer(ls)
		}
		return fullNode, err
	})


	stack.Start()

	fmt.Println("Started!")

	var ethereum *eth.Ethereum

	if err := stack.Service(&ethereum); err != nil {
		fmt.Println("Eth not running: ",err)
	} else {
		val := stack.GetService(&ethereum)
		api := val.(*eth.Ethereum).ApiBackend)
		api.GetReceipts()


	}
}
func main()  {

	/*
	handler := log.StreamHandler(os.Stdout, log.LogfmtFormat())
	log.Root().SetHandler(handler)
	log.Info("hahhah")
	get_database()
	return
	*/

	/*request_storage(
		"/Users/lagoon/Library/Ethereum/geth/chaindata",
		"fdaa1853593f060fa0fd2b35c8364daaebcd41a984cc4550858f7265a3d5df0d",
		"0x8d12A197cB00D4747a1fe03395095ce2A5CC6819",
		crypto.Keccak256Hash(
			common.Hex2BytesFixed("Bc69eCC478d4C1722A7151E43E2A78CCD0D5D5fB", 32),
			crypto.Keccak256Hash(
				common.Hex2BytesFixed("00", 32),
				common.Hex2BytesFixed("06", 32)).Bytes()).Bytes(),
	)
	*/


	request_proof(
		"/Users/lagoon/Library/Ethereum/geth/chaindata",
		"fdaa1853593f060fa0fd2b35c8364daaebcd41a984cc4550858f7265a3d5df0d",
		"0x8d12A197cB00D4747a1fe03395095ce2A5CC6819",
		"Bc69eCC478d4C1722A7151E43E2A78CCD0D5D5fB",
		"00",
		"06",
	)
	chaindb_dir := "/Users/lagoon/Library/Ethereum/geth/chaindata"
	block_number := uint64(4249995)
	//contract_address := "0x8d12A197cB00D4747a1fe03395095ce2A5CC6819"
	//user_address := "0xBc69eCC478d4C1722A7151E43E2A78CCD0D5D5fB"

	/*request_transaction(
		chaindb_dir,
		"0x93cbad7a82e339dc055a6d42be5fb04e2ba585a9181c9358fee3498581485fcc",
		block_number,
	)*/

	enc_header, header := request_header(
		chaindb_dir,
		block_number,
	)
	fmt.Println("header rlp:",enc_header)
	fmt.Println("header: ",header)
	/*
	enc_account_proof, account := request_account(
		chaindb_dir,
		header.Root.Hex(),
		contract_address,
		)
	fmt.Println("account proof rlp: ",enc_account_proof)
	fmt.Println("account :",account)




	enc_proof, value := request_proof(
		chaindb_dir,
		header.Root.Hex(),
		contract_address,
		user_address,
		"00",
		"06",
	)
	fmt.Println("proof:",enc_proof)
	fmt.Println("value: ",value)
	*/

}

package nodes

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/chenzhijie/go-web3"
	"github.com/chenzhijie/go-web3/eth"
	"github.com/chenzhijie/go-web3/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	eTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"go-venice/configs"
	"go-venice/pkg/utils"
	"io"
	"math/big"
	"net/http"
	"strings"
)

type EvmRPCMethod string
type EvmContractMethod string
type ContractAddress string

const (
	MethodEthGetBalance         EvmRPCMethod = "eth_getBalance"
	MethodEthSendRawTransaction EvmRPCMethod = "eth_sendRawTransaction"
)
const (
	BalanceOf      EvmContractMethod = "balanceOf"
	PendingRewards EvmContractMethod = "pendingRewards"
	Claim          EvmContractMethod = "claim"
	Approve        EvmContractMethod = "approve"
	Stake          EvmContractMethod = "stake"
)
const (
	TonkenAddress  ContractAddress = "0xacfe6019ed1a7dc6f7b508c02d1b04ec88cc21bf"
	StakingAddress ContractAddress = "0x321b7ff75154472B18EDb199033fF4D116F340Ff"
)

type Evm struct {
	url     string
	web3    *web3.Web3
	chainId *big.Int
}

func NewEvm(cfg *configs.EnvConfig) *Evm {
	w, err := web3.NewWeb3(cfg.RpcUrl)
	if err != nil {
		panic(err)
	}
	w.Eth.SetChainId(cfg.ChainId)
	return &Evm{cfg.RpcUrl, w, big.NewInt(cfg.ChainId)}
}

func (evm *Evm) GetTransaction(hash string) (*eTypes.Transaction, error) {
	tx, err := evm.web3.Eth.GetTransactionByHash(common.HexToHash(hash))
	if err != nil {
		return nil, err
	}
	fmt.Println(hash, " -> Transaction: ", tx)

	return tx, nil
}

func (evm *Evm) GetBlockNumber() (uint64, error) {
	blockNumber, err := evm.web3.Eth.GetBlockNumber()
	if err != nil {
		return 0, err
	}
	fmt.Println("Current block number: ", blockNumber)

	return blockNumber, nil
}
func (evm *Evm) GetBalance(ctx context.Context, address string) (*big.Int, error) {

	res, err := evm.call(ctx, MethodEthGetBalance, []string{address, "latest"})
	if err != nil {
		return nil, errors.Wrap(err, "GetBalance call fail")
	}
	balanceWei, err := utils.HexToBigInt(res.Result)
	if err != nil {
		return nil, errors.Wrap(err, "HexToBigInt fail")
	}

	fmt.Println(res.Result)
	return balanceWei, nil
}

func (evm *Evm) GetBalanceToken(address string) (*big.Int, error) {
	return evm.balanceOf(TonkenAddress, address)
}

func (evm *Evm) GetDelegated(address string) (*big.Int, error) {
	return evm.balanceOf(StakingAddress, address)
}

func (evm *Evm) GetReward(address string) (*big.Int, error) {
	abiString := `[
	{
	  "inputs": [
		{
		  "internalType": "address",
		  "name": "_user",
		  "type": "address"
		}
	  ],
	  "name": "pendingRewards",
	  "outputs": [
		{
		  "internalType": "uint256",
		  "name": "",
		  "type": "uint256"
		}
	  ],
	  "stateMutability": "view",
	  "type": "function"
	}
]`
	contract, err := evm.web3.Eth.NewContract(abiString, string(StakingAddress))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	call, err := contract.Call(string(PendingRewards), common.HexToAddress(address))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	fmt.Println(call)
	return call.(*big.Int), nil
}

func (evm *Evm) CreateTransferTransaction(fromAddress, toAddress, ethAmount string) (string, error) {
	return evm.createRawTransaction(fromAddress, toAddress, ethAmount, nil)

}
func (evm *Evm) CreateClaimTransaction(fromAddress string) (string, error) {
	abiString := `[
	{
	  "inputs": [],
	  "name": "claim",
	  "outputs": [],
	  "stateMutability": "nonpayable",
	  "type": "function"
	}
]`
	contract, err := evm.web3.Eth.NewContract(abiString, string(StakingAddress))
	if err != nil {
		return "", errors.Wrap(err, "NewContract error")
	}
	abi, err := contract.EncodeABI(string(Claim))
	if err != nil {
		return "", err
	}

	fmt.Println("abi: ", hexutil.Encode(abi))

	transaction, err := evm.createRawTransaction(fromAddress, string(StakingAddress), "0", abi)
	if err != nil {
		return "", errors.Wrap(err, "CreateRawTransaction error")
	}

	return transaction, nil
}

func (evm *Evm) CreateApproveTransaction(fromAddress string, ethAmount string) (string, error) {
	abiString := `[
	{
	  "inputs": [
		{
		  "internalType": "address",
		  "name": "spender",
		  "type": "address"
		},
		{
		  "internalType": "uint256",
		  "name": "amount",
		  "type": "uint256"
		}
	  ],
	  "name": "approve",
	  "outputs": [
		{
		  "internalType": "bool",
		  "name": "",
		  "type": "bool"
		}
	  ],
	  "stateMutability": "nonpayable",
	  "type": "function"
	}
]`
	contract, err := evm.web3.Eth.NewContract(abiString, string(TonkenAddress))
	if err != nil {
		return "", errors.Wrap(err, "NewContract error")
	}
	abi, err := contract.EncodeABI(string(Approve), common.HexToAddress(string(StakingAddress)), evm.web3.Utils.ToWei(ethAmount))
	if err != nil {
		return "", err
	}

	fmt.Println("abi: ", hexutil.Encode(abi))

	transaction, err := evm.createRawTransaction(fromAddress, string(TonkenAddress), "0", abi)
	if err != nil {
		return "", errors.Wrap(err, "CreateRawTransaction error")
	}

	return transaction, nil
}

func (evm *Evm) CreateStakeTransaction(fromAddress string, ethAmount string) (string, error) {
	abiString := `[
	{
	  "inputs": [
		{
		  "internalType": "address",
		  "name": "recipient",
		  "type": "address"
		},
		{
		  "internalType": "uint256",
		  "name": "amount",
		  "type": "uint256"
		}
	  ],
	  "name": "stake",
	  "outputs": [],
	  "stateMutability": "nonpayable",
	  "type": "function"
	}
]`
	contract, err := evm.web3.Eth.NewContract(abiString, string(StakingAddress))
	if err != nil {
		return "", errors.Wrap(err, "NewContract error")
	}
	abi, err := contract.EncodeABI(string(Stake), common.HexToAddress(fromAddress), evm.web3.Utils.ToWei(ethAmount))
	if err != nil {
		return "", err
	}

	fmt.Println("abi: ", hexutil.Encode(abi))

	transaction, err := evm.createRawTransaction(fromAddress, string(StakingAddress), "0", abi)
	if err != nil {
		return "", errors.Wrap(err, "CreateRawTransaction error")
	}

	return transaction, nil
}

func (evm *Evm) createRawTransaction(fromAddress, toAddress, ethAmount string, input []byte) (string, error) {
	nonce, err := evm.getNonce(fromAddress)
	if err != nil {
		return "", err
	}
	estimateFee, err := evm.GetEstimateFee()
	if err != nil {
		return "", err
	}

	gasLimit, err := evm.GetEstimateGas(fromAddress, toAddress, input)
	if err != nil {
		return "", err
	}
	fmt.Println("MaxPriorityFeePerGas: ", estimateFee.MaxPriorityFeePerGas)
	fmt.Println("BaseFee: ", estimateFee.BaseFee)
	fmt.Println("gasLimit: ", gasLimit)
	to := common.HexToAddress(toAddress)
	dynamicFeeTx := &eTypes.DynamicFeeTx{
		ChainID:   evm.chainId,
		Nonce:     nonce,
		GasTipCap: estimateFee.MaxPriorityFeePerGas,
		GasFeeCap: estimateFee.MaxFeePerGas,
		Gas:       gasLimit,
		To:        &to,
		Value:     evm.web3.Utils.ToWei(ethAmount),
		Data:      input,
	}
	encodedTx, err := rlp.EncodeToBytes(dynamicFeeTx)
	if err != nil {
		return "", err
	}
	hexTx := "0x02" + common.Bytes2Hex(encodedTx)
	return hexTx, nil
}

func (evm *Evm) SingRawTransaction(unsigned string, privateKeyHex string) (string, error) {
	raw, err := hex.DecodeString(strings.TrimPrefix(unsigned, "0x")) // "0x" 제거
	if err != nil {
		return "", errors.Wrap(err, "hex decode 실패")
	}
	if len(raw) == 0 || raw[0] != 0x02 {
		return "", errors.New("not a dynamic fee tx (type 0x02)")
	}

	var tx eTypes.DynamicFeeTx
	if err := rlp.DecodeBytes(raw[1:], &tx); err != nil {
		return "", errors.Wrap(err, "DecodeBytes 실패")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", errors.Wrap(err, "HexToECDSA 실패")
	}
	signedTx, err := eTypes.SignNewTx(privateKey, eTypes.LatestSignerForChainID(evm.chainId), &tx)
	if err != nil {
		return "", errors.Wrap(err, "SignNewTx 실패")
	}
	txData, err := signedTx.MarshalBinary()
	return hexutil.Encode(txData), nil
}

func (evm *Evm) Broadcast(ctx context.Context, signedTx string) (string, error) {

	resp, err := evm.call(ctx, MethodEthSendRawTransaction, []string{signedTx})
	if err != nil {
		return "", errors.Wrap(err, "MethodEthSendRawTransaction 실패")
	}

	return resp.Result, nil
}

func (evm *Evm) getNonce(address string) (uint64, error) {
	nonce, err := evm.web3.Eth.GetNonce(common.HexToAddress(address), nil)
	if err != nil {
		return 0, err
	}
	fmt.Println("Latest nonce: ", nonce)

	return nonce, nil
}

func (evm *Evm) GetGasPrice() (uint64, error) {
	gasPrice, err := evm.web3.Eth.GasPrice()
	if err != nil {
		return 0, err
	}
	fmt.Println("gasPrice: ", gasPrice)

	return gasPrice, nil
}

func (evm *Evm) GetEstimateFee() (*eth.EstimateFee, error) {
	fee, err := evm.web3.Eth.EstimateFee()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	fmt.Println("EstimateFee: ", fee)

	return fee, nil
}

func (evm *Evm) GetEstimateGas(from, to string, data types.CallMsgData) (uint64, error) {
	fee, err := evm.web3.Eth.EstimateGas(&types.CallMsg{
		From: common.HexToAddress(from), To: common.HexToAddress(to), Data: data, Gas: nil, GasPrice: nil, Value: nil,
	})
	if err != nil {
		return 0, err
	}
	fmt.Println("EstimateFee: ", fee)

	return fee, nil
}

func (evm *Evm) GetEstimatePriorityFee() (*big.Int, error) {
	blockNumber, err := evm.GetBlockNumber()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	gasPrice, err := evm.web3.Eth.EstimatePriorityFee(3, big.NewInt(int64(blockNumber)), []float64{50})
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	fmt.Println("EstimatePriorityFee: ", gasPrice)

	return gasPrice, nil
}

func (evm *Evm) GetFeeHistory() (*types.FeeHistory, error) {
	blockNumber, err := evm.GetBlockNumber()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	feeHistory, err := evm.web3.Eth.FeeHistory(1, big.NewInt(int64(blockNumber)), []float64{50})
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	fmt.Println("feeHistory: ", feeHistory)

	return feeHistory, nil
}

func (evm *Evm) balanceOf(contractAddress ContractAddress, address string) (*big.Int, error) {
	abiString := `[
	{
	  "inputs": [
		{
		  "internalType": "address",
		  "name": "account",
		  "type": "address"
		}
	  ],
	  "name": "balanceOf",
	  "outputs": [
		{
		  "internalType": "uint256",
		  "name": "",
		  "type": "uint256"
		}
	  ],
	  "stateMutability": "view",
	  "type": "function"
	}
]`
	contract, err := evm.web3.Eth.NewContract(abiString, string(contractAddress))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	call, err := contract.Call(string(BalanceOf), common.HexToAddress(address))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	fmt.Println(call)
	return call.(*big.Int), nil
}

func createPayload(method EvmRPCMethod, params []string) (io.Reader, error) {
	req := EvmRPCRequest{
		JsonRpc: "2.0",
		Method:  string(method),
		Id:      1,
		Params:  params,
	}

	bytesRequest, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return bytes.NewReader(bytesRequest), nil
}

func (evm *Evm) call(ctx context.Context, method EvmRPCMethod, params []string) (*EvmRPCResponse, error) {
	payload, err := createPayload(method, params)
	if err != nil {
		fmt.Println(err)
		return nil, errors.Wrap(err, "")
	}
	req, err := http.NewRequestWithContext(ctx, "POST", evm.url, payload)
	if err != nil {
		fmt.Println(err)
		return nil, errors.Wrap(err, "")
	}
	req.Header.Add("Content-Type", "application/json")
	defer req.Body.Close()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, errors.Wrap(err, "DefaultClient fail")
	}
	fmt.Println(res.Body)
	byteBody, err := io.ReadAll(res.Body)

	resDto := EvmRPCResponse{}
	if err := json.Unmarshal(byteBody, &resDto); err != nil {
		return nil, errors.Wrap(err, "")
	}
	if resDto.Error != nil {
		return nil, fmt.Errorf("rpc error: %s (code %d)", resDto.Error.Message, resDto.Error.Code)
	}
	return &resDto, nil
}

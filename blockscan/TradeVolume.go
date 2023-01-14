package blockscan

import (
        "context"
        // "fmt"
        // "github.com/HydroProtocol/ethereum-watcher/plugin"
        ethereum_watcher "github.com/HydroProtocol/ethereum-watcher"
        "github.com/sirupsen/logrus"

        // "github.com/shopspring/decimal"
        "github.com/HydroProtocol/ethereum-watcher/blockchain"
        "github.com/HydroProtocol/ethereum-watcher/plugin"
        // "github.com/ethereum/go-ethereum/ethdb/leveldb"
)

const (
        CONTRACT = "0x55d398326f99059fF775485246999027B3197955"

        // 收款地址 0x34A8417d3747eB3c80d9CC94522C4eb89c95D967
        // BEE-USDT 0xe94a77cbcc0590b0738c48b2594c08e456fed89d
        BEE_USDT_PAIR = "0x000000000000000000000000e94a77cbcc0590b0738c48b2594c08e456fed89d"
)

var ERC20_transfer = []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}

func main() {

        api := "https://fragrant-dark-pine.bsc.quiknode.pro/977be32e92f335ddc7095549c7264aeedb95ba23/"
        startBlockNum := 17926307
        contractAdx := CONTRACT

        handler := func(from, to int, receiptLogs []blockchain.IReceiptLog, isUpToHighestBlock bool) error {

                for _, log := range receiptLogs {

                        Topics := log.GetTopics()
                        from_addr := Topics[1]
                        to_addr := Topics[2]
                        token := log.GetAddress()
                        weiamount, _ := plugin.HexToDecimal(log.GetData())

                        logrus.Infof("订单成交: %s", order_id)

                }

                return nil
        }

        receiptLogWatcher := ethereum_watcher.NewReceiptLogWatcher(
                context.Background(),
                api,
                startBlockNum,
                contractAdx,
                ERC20_transfer,
                handler,
                ethereum_watcher.ReceiptLogWatcherConfig{
                        StepSizeForBigLag:               5,
                        IntervalForPollingNewBlockInSec: 5,
                        RPCMaxRetry:                     3,
                        ReturnForBlockWithNoReceiptLog:  true,
                },
        )

        receiptLogWatcher.Run()

}

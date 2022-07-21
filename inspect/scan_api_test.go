package inspect

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"
)

const (
	testScanAPIOldestSupportedBlock = uint64(0)
)

var testScanEnv struct {
	ScanAPI string `envconfig:"scan_api" default:"https://cloudflare-eth.com"`
}

func init() {
	envconfig.MustProcess("test", &testScanEnv)
}

func TestScanAPIInspection(t *testing.T) {
	r := require.New(t)

	inspector := &ScanAPIInspector{}
	results, err := inspector.Inspect(context.Background(), InspectionConfig{
		ScanAPIURL:  testScanEnv.ScanAPI,
		BlockNumber: testScanAPIOldestSupportedBlock,
	})
	r.NoError(err)

	r.Equal(map[string]float64{
		IndicatorScanAPIAccessible:     ResultSuccess,
		IndicatorScanAPIChainID:        float64(1),
		IndicatorScanAPIModuleWeb3:     ResultSuccess,
		IndicatorScanAPIModuleEth:      ResultSuccess,
		IndicatorScanAPIModuleNet:      ResultSuccess,
		IndicatorScanAPIHistorySupport: VeryOldBlockNumber,
	}, results.Indicators)

	r.Equal(map[string]string{
		MetadataScanAPIBlockByNumberHash: "3abe2f22edf2b463cbc13343a947be9ebbf8c16c2b50b2b90e10a199a2344f65",
	}, results.Metadata)
}

func Test_findOldestSupportedBlock(t *testing.T) {
	r := require.New(t)

	cli, err := ethclient.Dial(testScanEnv.ScanAPI)
	r.NoError(err)

	ctx := context.Background()
	latestBlockNum, err := cli.BlockNumber(ctx)
	r.NoError(err)

	result := findOldestSupportedBlock(context.Background(), cli, 0, latestBlockNum)
	r.Equal(testScanAPIOldestSupportedBlock, result)
}

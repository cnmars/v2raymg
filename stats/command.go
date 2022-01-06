package stats

import (
	"context"
	"fmt"
	"regexp"

	"github.com/lureiny/v2raymg/config"
	"github.com/v2fly/v2ray-core/v4/app/stats/command"
	"google.golang.org/grpc"
)

// MyStat 集成了用户uplink和downlink的
type MyStat struct {
	Name     string
	Type     string
	Downlink int64
	Uplink   int64
}

var regexCompile = regexp.MustCompile(`(user|inbound|outbound)>>>(\S+)>>>traffic>>>(downlink|uplink)`)

func queryStats(con command.StatsServiceClient, req *command.QueryStatsRequest) (*command.QueryStatsResponse, error) {
	resp, err := con.QueryStats(context.Background(), req)
	if err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

func QueryStats(host string, port int) (*command.QueryStatsResponse, error) {
	// 创建grpc client
	cmdConn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		config.Error.Fatal(err)
	}

	statClient := command.NewStatsServiceClient(cmdConn)

	// query 参数
	queryStatsReq := command.QueryStatsRequest{
		Pattern: "",
		Reset_:  false,
	}
	resp, err := statClient.QueryStats(context.Background(), &queryStatsReq)

	if err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

func queryAllStats(con command.StatsServiceClient, req *command.QueryStatsRequest) (*map[string]*MyStat, error) {
	stats, err := queryStats(con, req)
	if err != nil {
		return nil, err
	}
	result := make(map[string]*MyStat)
	for _, stat := range stats.GetStat() {
		reResult := regexCompile.FindStringSubmatch(stat.GetName())
		if _, ok := result[reResult[2]]; !ok {
			result[reResult[2]] = &MyStat{
				Name: reResult[2],
				Type: reResult[1],
			}
		}
		// 填充数据流量
		switch reResult[3] {
		case "downlink":
			result[reResult[2]].Downlink = stat.GetValue()
		case "uplink":
			result[reResult[2]].Uplink = stat.GetValue()
		}
	}
	return &result, nil
}

func QueryAllStats(host string, port int) (*map[string]*MyStat, error) {
	// 创建grpc client
	cmdConn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		config.Error.Fatal(err)
	}

	statClient := command.NewStatsServiceClient(cmdConn)

	// query 参数
	queryStatsReq := command.QueryStatsRequest{
		Pattern: "",
		Reset_:  false,
	}
	return queryAllStats(statClient, &queryStatsReq)

}

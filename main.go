package main

import (
	"bytes"
	"context"
	"fmt"
	"encoding/json"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
        "github.com/sensu/sensu-go/types"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

type Config struct {
	sensu.PluginConfig
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "check-elasticsearch-health",
			Short:    "A very basic Elasticsearch health check",
			Keyspace: "sensu.io/plugins/check-elasticsearch-health/config",
		},
	}

	options = []*sensu.PluginConfigOption{
	}
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
	    fmt.Printf("Critical: %v\n", err)
	    return sensu.CheckStateCritical, nil
	}	

	req := esapi.ClusterHealthRequest{}
	res, err := req.Do(context.Background(), es)
	if err != nil {
	    fmt.Printf("Critical: %v\n", err)
	    return sensu.CheckStateCritical, nil
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(res.Body)
	err != nil {
	    fmt.Printf("Critical: %v\n", err)
	    return sensu.CheckStateCritical, nil
	}
        result := buf.String()

        var final map[string]interface{}
	err = json.Unmarshal([]byte(result), &final)
	if err != nil {
	    fmt.Printf("Critical: %v\n", err)
	    return sensu.CheckStateCritical, nil
	}
 
	if final["status"] == "green" {
	   fmt.Printf("%s OK: cluster status is Green.\n", plugin.PluginConfig.Name)
	   return sensu.CheckStateOK, nil
	} else if final["status"] == "yellow" {
	   fmt.Printf("%s Warning: cluster status is Yellow.\n", plugin.PluginConfig.Name)
	   return sensu.CheckStateWarning, nil
	} else if final["status"] == "red" {
	    fmt.Printf("%s Critical: cluster status is Red.\n", plugin.PluginConfig.Name)
	    return sensu.CheckStateCritical, nil
	} else {
	    fmt.Printf("%s UNKNOWN: cluster has no status!\n", plugin.PluginConfig.Name)
	    return sensu.CheckStateUnknown, nil
	}

}

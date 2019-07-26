// package main
package main

import (
	"context"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"reflect"

	"github.com/golang/protobuf/proto"    
)

// Define global variable releases to hold name of benchmarks and their keys.
var releases = map[string]map[string]interface{}{
	"master": map[string]interface{}{
		"config_count": "6086318920564736",
		"cpu":          "4852237377470464",
		"grpc":         "4866368222527488",
		"memory":       "6140819840958464",
		"requests":     "6564959337054208",
	},
}

// Define global variable metrics to list metrics for each benchmark in releases.
var metrics = map[string]map[string]string{
	"config_count": map[string]string{
		"mixer_config_rule_config_match_error_count":    "m1",
		"mixer_config_unsatisfied_action_handler_count": "m2",
		"mixer_config_instance_config_count":            "m3",
		"mixer_config_rule_config_error_count":          "m4",
		"mixer_config_rule_config_count":                "m5",
		"mixer_config_attribute_count":                  "m6",
	},
	"cpu": map[string]string{
		"cpu_mili_policy_proxy":    "m1",
		"cpu_mili_telemetry_proxy": "m2",
		"cpu_mili_pilot_discovery": "m3",
		"cpu_mili_telemetry_mixer": "m4",
		"cpu_mili_policy_mixer":    "m5",
	},
	"grpc": map[string]string{
		"grpc_server_handled_total_4xx": "m1",
		"grpc_server_handled_total":     "m2",
		"grpc_server_handled_total_5xx": "m3",
	},
	"memory": map[string]string{
		"mem_MB_policy_mixer":    "m1",
		"mem_MB_telemetry_proxy": "m2",
		"mem_MB_policy_proxy":    "m3",
		"mem_MB_pilot_discovery": "m4",
		"mem_MB_telemetry_mixer": "m5",
	},
	"requests": map[string]string{
		"istio_requests_total_503": "m1",
		"istio_requests_total_504": "m2",
		"istio_requests_total_404": "m3",
	},
}

type pointValues struct {
	time        float64
	samplePoint map[string]float64
}

func equals(p1, p2 pointValues) bool {
	if p1.time != p2.time {
		return false
	}
	eq := reflect.DeepEqual(p1.samplePoint, p2.samplePoint)
	if eq {
		return true
	}
	return false
}

func contains(p pointValues, pValues []pointValues) bool {
	for _, pValue := range(pValues) {
		if equals(pValue, p) {
			return true
		}
	}
	return false
}

func sum(input []interface{}) float64 {
	var sum float64
	for _, num := range input {
		checkInt, err := num.(float64)
		if err {
			sum += checkInt
		}
	}
	return sum
}

// Join branch map from current iteration to a total branch map.
func add(totalBranchMap map[string]map[interface{}][]pointValues, branchMap map[string]map[interface{}][]pointValues) map[string]map[interface{}][]pointValues {
	for branchName, branchContent := range(branchMap) {
		if totalBranchMap[branchName] == nil {
			totalBranchMap[branchName] = branchContent
		} else {
			totalBranchContent := totalBranchMap[branchName]
			for benchKey, pValues := range(branchContent) {
				if totalBranchContent[benchKey] == nil {
					totalBranchContent[benchKey] = pValues
				} else {
					totalBenchPoints := totalBranchContent[benchKey]
					for _, pValue := range(pValues) {
						if !contains(pValue, totalBenchPoints) {
							totalBenchPoints = append(totalBenchPoints, pValue)
						}
					}
				}
			}
		}
	}
	return totalBranchMap
}

// Process the metric value of each element in json-formatted output.
// Convert it to float64 based on its original type.
func processMetricContent(originalValueForMetric interface{}) float64 {
	var floatValue float64
	stringValue, err := originalValueForMetric.(string)
	if err {
		if strings.Compare(stringValue, "NaN") == 0 {
			fmt.Println("value is NaN computed value is 0.")
			floatValue = 0
		} else if stringOfIntValue, err := strconv.Atoi(stringValue); err == nil {
			fmt.Printf("value is string of number, computed value is %d.", stringOfIntValue)
			floatValue = float64(stringOfIntValue)
		}
	} else if sliceValue, err := originalValueForMetric.([]interface{}); err {
		sumOfSlice := sum(sliceValue)
		fmt.Printf("value is list of numbers, computed value is %d.", sumOfSlice)
		floatValue = sumOfSlice
	} else if intValue, err := originalValueForMetric.(int); err {
		fmt.Printf("value is integer, computed value is %d.", intValue)
		floatValue = float64(intValue)
	}
	return floatValue
}

// Store process map with benchmarks to quickstore mako.
func storeValueToMako(branchMap map[string]map[interface{}][]pointValues) {
  // Temporary code: run docker command to start docker container.
  // `docker.sh` contains code `docker run --rm --name=mako-storage -v ~/.config/gcloud/application_default_credentials.json:/root/adc.json -e "GOOGLE_APPLICATION_CREDENTIALS=/root/adc.json" -p 9813:9813 us.gcr.io/makoperf/mako-microservice:latest`
  cmd := exec.Command("/bin/sh", "docker.sh")
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Error running docker command: %v", err)
	}
	addr := "localhost:9813"
    fmt.Printf("Connecting to microservice at %s", addr)
    ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)
	for _, benchMaps := range branchMap {
		for benchKey, samplePoints := range benchMaps {

		    // Leave closeq function blank to keep the docker container running for other benchmark keys.
		    q, _, err := quickstore.NewAtAddress(ctxWithTimeout, &qpb.QuickstoreInput{BenchmarkKey: proto.String(benchKey.(string))}, addr)
		    if err != nil {
		            fmt.Printf("Uh oh...\n")
		            log.Fatalf("failed NewAtAddress: %v", err)
		    }

			for _, pointPair := range samplePoints {
				err := q.AddSamplePoint(pointPair.time, pointPair.samplePoint)
				if err != nil {
					log.Fatalf("AddSamplePoint err: %v", err)
				}
			}
			out, err := q.Store()
			if err != nil {
				log.Fatalf("Store() err: %v", err)
			}
			fmt.Printf("View chart: %s", out.GetRunChartLink())
		}
	}
	
	cancel()
  	// Temporary command: close docker container.
	// `dockerstop.sh` contains code `docker stop`
	cmd = exec.Command("/bin/sh", "dockerstop.sh")
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Error running docker command: %v", err)
	}
	
}

// Check if output in json format contains metrics defined for the benchmarks.
// Gather all points for each benchmark key.
func processMarshaledOutput(marshaledOutput map[string]interface{}) map[string]map[interface{}][]pointValues {
	branchMap := map[string]map[interface{}][]pointValues{}
	for branch := range releases {
		eachBenchnameMap := map[interface{}][]pointValues{}
		for benchname, benchkey := range releases[branch] {
			benchMap := []pointValues{}
			metricMap := metrics[benchname]
			for metricKeyName, originalValueForMetric := range marshaledOutput {
				ts := float64(time.Now().UnixNano() / time.Millisecond.Nanoseconds())
				if metricMap[metricKeyName] != "" {
					metricKeyIndex := metricMap[metricKeyName]
					floatValue := processMetricContent(originalValueForMetric)
					timeAndPoint := pointValues{time: ts, samplePoint: map[string]float64{metricKeyIndex: floatValue}}
					benchMap = append(benchMap, timeAndPoint)
				}
				ts++
			}
			if len(benchMap) != 0 {
				eachBenchnameMap[benchkey] = benchMap
			}
		}
		if len(eachBenchnameMap) != 0 {
			branchMap[branch] = eachBenchnameMap
		}
	}
	fmt.Println(branchMap)
	return branchMap
}

// Run prom.py and process its output.
func getMetricResultFromPython(timeToRunPerfTests int) {
  // Run prom.py from istio.io.
	runPromPath := "runProm.sh"
	fmt.Println("Get metric %d time(s).", timeToRunPerfTests)
	totalBranchMap := map[string]map[interface{}][]pointValues{}
	for i := 0; i < timeToRunPerfTests; i++ {
		fmt.Println("Running for %d min", i)

		// Call bash script to run prom.py and read stdout from command line.
		cmd := exec.Command("/bin/sh", runPromPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Error reading output from prom: %v", err)
		}

		// Process stdout into json format as interface{}.
		var v interface{}
		err = json.Unmarshal([]byte(out.String()), &v)
		if err != nil {
			log.Fatalf("Error unmarshalling json: %v", err)
		}
		marshaledOutput := v.(map[string]interface{})

		// Process the output to match benchmark keys with point values.
		branchMap := processMarshaledOutput(marshaledOutput)

		// Add branch map for each runs to a total branch map to store points together later.
		totalBranchMap = add(totalBranchMap, branchMap)
	}
	storeValueToMako(totalBranchMap)
}

func main() {
	var timeToRunPerfTests int
	flag.IntVar(&timeToRunPerfTests, "TIME_TO_RUN_PERF_TESTS", 2, "Input time to run perf tests")

	// Need to call flag.Parse not only to parse TIME_TO_RUN_PERF_TESTS but also to avoid logging errors from glog.
	flag.Parse()
	if timeToRunPerfTests == 0 {
		log.Fatalf("Value of TIME_TO_RUN_PERF_TESTS must be greater than 0.")
	}

	// Call prom.py from istio.io/tools.
	getMetricResultFromPython(timeToRunPerfTests)
}

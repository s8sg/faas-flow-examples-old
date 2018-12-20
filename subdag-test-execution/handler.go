package function

import (
	faasflow "github.com/s8sg/faas-flow"
	consulStateStore "github.com/s8sg/faas-flow-consul-statestore"
	sdk "github.com/s8sg/faas-flow/sdk"
	"os"
	"time"
)

func debug(node string) sdk.Modifier {
	return func(data []byte) ([]byte, error) {
		if len(data) == 0 {
			// Get initial time
			data = []byte("Start: " + time.Now().String())
		}
		data = []byte(node + "( " + string(data) + " )")
		return data, nil
	}
}

func aggregator(nodes map[string][]byte) ([]byte, error) {
	data := ""
	for _, value := range nodes {
		if len(data) == 0 {
			data = string(value)
		}
		data = string(value) + " + " + data
	}
	return []byte(data), nil
}

// Define provide definiton of the workflow
func Define(flow *faasflow.Workflow, context *faasflow.Context) (err error) {

	// Create the dag
	maindag := faasflow.CreateDag()

	maindag.AddModifier("n1", debug("n1"))
	maindag.AddModifier("n2", debug("n2"))
	maindag.AddModifier("n3", debug("n3"))
	maindag.AddModifier("n4", debug("n4"))

	maindag.AddVertex("n5", faasflow.Aggregator(aggregator))
	maindag.AddModifier("n5", debug("n5"))

	maindag.AddModifier("n6", debug("n6"))

	maindag.AddVertex("n7", faasflow.Aggregator(aggregator))
	maindag.AddModifier("n7", debug("n7"))

	maindag.AddModifier("n8", debug("n8"))
	maindag.AddModifier("n9", debug("n9"))
	maindag.AddModifier("n10", debug("n10"))

	maindag.AddVertex("n11", faasflow.Aggregator(aggregator))
	maindag.AddModifier("n11", debug("n11"))

	subdag := faasflow.CreateDag()

	subdag.AddModifier("n13", debug("n13"))
	subdag.AddModifier("n14", debug("n14"))

	superSubDag := faasflow.CreateDag()
	superSubDag.AddModifier("ss1", debug("ss1"))
	superSubDag.AddModifier("ss2", debug("ss2"))
	superSubDag.AddModifier("ss3", debug("ss3"))
	superSubDag.AddModifier("ss4", debug("ss4"))
	superSubDag.AddVertex("ss5", faasflow.Aggregator(aggregator))
	superSubDag.AddModifier("ss5", debug("ss5"))

	superSubDag.AddEdge("ss1", "ss2", faasflow.Execution)
	superSubDag.AddEdge("ss2", "ss3", faasflow.Execution)
	superSubDag.AddEdge("ss3", "ss5", faasflow.Execution)
	superSubDag.AddEdge("ss1", "ss4", faasflow.Execution)
	superSubDag.AddEdge("ss4", "ss5", faasflow.Execution)

	subdag.AddSubDag("n15", superSubDag)

	subdag.AddVertex("n16", faasflow.Aggregator(aggregator))
	subdag.AddModifier("n16", debug("n16"))
	subdag.AddEdge("n13", "n14", faasflow.Execution)
	subdag.AddEdge("n13", "n15", faasflow.Execution)
	subdag.AddEdge("n14", "n16", faasflow.Execution)
	subdag.AddEdge("n15", "n16", faasflow.Execution)

	maindag.AddSubDag("n12", subdag)

	maindag.AddVertex("n17", faasflow.Aggregator(aggregator))
	maindag.AddModifier("n17", debug("n17"))
	maindag.AddEdge("n1", "n2", faasflow.Execution)
	maindag.AddEdge("n2", "n3", faasflow.Execution)
	maindag.AddEdge("n2", "n4", faasflow.Execution)
	maindag.AddEdge("n2", "n6", faasflow.Execution)
	maindag.AddEdge("n3", "n5", faasflow.Execution)
	maindag.AddEdge("n4", "n5", faasflow.Execution)
	maindag.AddEdge("n5", "n7", faasflow.Execution)
	maindag.AddEdge("n6", "n7", faasflow.Execution)
	maindag.AddEdge("n7", "n17", faasflow.Execution)

	maindag.AddEdge("n1", "n8", faasflow.Execution)
	maindag.AddEdge("n8", "n9", faasflow.Execution)
	maindag.AddEdge("n8", "n10", faasflow.Execution)
	maindag.AddEdge("n9", "n11", faasflow.Execution)
	maindag.AddEdge("n10", "n11", faasflow.Execution)
	maindag.AddEdge("n11", "n17", faasflow.Execution)

	maindag.AddEdge("n1", "n12", faasflow.Execution)
	maindag.AddEdge("n12", "n17", faasflow.Execution)

	flow.ExecuteDag(maindag)

	return
}

// DefineStateStore provides the override of the default StateStore
func DefineStateStore() (faasflow.StateStore, error) {
	consulss, err := consulStateStore.GetConsulStateStore(os.Getenv("consul_url"), os.Getenv("consul_dc"))
	if err != nil {
		return nil, err
	}
	return consulss, nil
}

// ProvideDataStore provides the override of the default DataStore
func DefineDataStore() (faasflow.DataStore, error) {
	return nil, nil
}

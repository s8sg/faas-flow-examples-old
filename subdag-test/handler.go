package function

import (
	faasflow "github.com/s8sg/faas-flow"
	consulStateStore "github.com/s8sg/faas-flow-consul-statestore"
	minioDataStore "github.com/s8sg/faas-flow-minio-datastore"
	sdk "github.com/s8sg/faas-flow/sdk"
	"log"
	"os"
	"time"
)

func debug(node string) sdk.Modifier {
	return func(data []byte) ([]byte, error) {
		if len(data) == 0 {
			// Get initial time
			t := time.Now()
			data = []byte(t.Format("15:04:05"))
		}
		data = []byte(node + "( " + string(data) + " )")
		log.Printf("Executing Node %s", node)
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

	superSubDag1 := faasflow.CreateDag()
	superSubDag1.AddModifier("ss1", debug("ss1"))
	superSubDag1.AddModifier("ss2", debug("ss2"))
	superSubDag1.AddModifier("ss3", debug("ss3"))
	superSubDag1.AddModifier("ss4", debug("ss4"))
	superSubDag1.AddVertex("ss5", faasflow.Aggregator(aggregator))
	superSubDag1.AddModifier("ss5", debug("ss5"))
	superSubDag1.AddEdge("ss1", "ss2")
	superSubDag1.AddEdge("ss2", "ss3")
	superSubDag1.AddEdge("ss3", "ss5")
	superSubDag1.AddEdge("ss1", "ss4")
	superSubDag1.AddEdge("ss4", "ss5")

	superSubDag2 := faasflow.CreateDag()
	superSubDag2.AddModifier("dd1", debug("dd1"))
	superSubDag2.AddModifier("dd2", debug("dd2"))
	superSubDag2.AddModifier("dd3", debug("dd3"))
	superSubDag2.AddModifier("dd4", debug("dd4"))
	superSubDag2.AddVertex("dd5", faasflow.Aggregator(aggregator))
	superSubDag2.AddModifier("dd5", debug("dd5"))
	superSubDag2.AddEdge("dd1", "dd2")
	superSubDag2.AddEdge("dd2", "dd3")
	superSubDag2.AddEdge("dd3", "dd5")
	superSubDag2.AddEdge("dd1", "dd4")
	superSubDag2.AddEdge("dd4", "dd5")

	// Normal Graph
	// subdag.AddSubDag("n15", superSubDag1)

	// Foreach Graph
	// subdag.AddForEachDag("n15", superSubDag1, faasflow.ForEach(func(data []byte) map[string][]byte {
	// 	return map[string][]byte{"c1": data, "c2": data, "c3": data}
	// }, aggregator))

	// Condition Graph
	subdag.AddConditionalDags("n15", map[string]*faasflow.DagFlow{"c1": superSubDag1, "c2": superSubDag2},
		faasflow.Condition(func(data []byte) []string {
			return []string{"c1", "c2"}
		}, aggregator))

	subdag.AddVertex("n16", faasflow.Aggregator(aggregator))
	subdag.AddModifier("n16", debug("n16"))
	subdag.AddEdge("n13", "n14")
	subdag.AddEdge("n13", "n15")
	subdag.AddEdge("n14", "n16")
	subdag.AddEdge("n15", "n16")
	maindag.AddSubDag("n12", subdag)

	maindag.AddVertex("n17", faasflow.Aggregator(aggregator))
	maindag.AddModifier("n17", debug("n17"))

	maindag.AddEdge("n1", "n2")
	maindag.AddEdge("n2", "n3")
	maindag.AddEdge("n2", "n4")
	maindag.AddEdge("n2", "n6")
	maindag.AddEdge("n3", "n5")
	maindag.AddEdge("n4", "n5")
	maindag.AddEdge("n5", "n7")
	maindag.AddEdge("n6", "n7")
	maindag.AddEdge("n7", "n17")

	maindag.AddEdge("n1", "n8")
	maindag.AddEdge("n8", "n9")
	maindag.AddEdge("n8", "n10")
	maindag.AddEdge("n9", "n11")
	maindag.AddEdge("n10", "n11")
	maindag.AddEdge("n11", "n17")

	maindag.AddEdge("n1", "n12")
	maindag.AddEdge("n12", "n17")

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
	// initialize minio DataStore
	miniods, err := minioDataStore.InitFromEnv()
	if err != nil {
		return nil, err
	}

	return miniods, nil
}

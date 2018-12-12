package function

import (
	faasflow "github.com/s8sg/faas-flow"
	consulStateStore "github.com/s8sg/faas-flow-consul-statestore"
	minioDataStore "github.com/s8sg/faas-flow-minio-datastore"
	sdk "github.com/s8sg/faas-flow/sdk"
	"os"
)

func debug(node string) sdk.Modifier {
	return func(data []byte) ([]byte, error) {
		data = []byte(node + " (" + string(data) + ") ")
		return data, nil
	}
}

func serializer(nodes map[string][]byte) ([]byte, error) {
	data := []byte("")
	for _, value := range nodes {
		data = []byte(string(data) + " + " + string(value))
	}
	return data, nil
}

// Define provide definiton of the workflow
func Define(flow *faasflow.Workflow, context *faasflow.Context) (err error) {

	// Create the dag
	maindag := faasflow.CreateDag()

	maindag.AddModifier("n1", debug("n1"))
	maindag.AddModifier("n2", debug("n2"))
	maindag.AddModifier("n3", debug("n3"))
	maindag.AddModifier("n4", debug("n4"))

	maindag.AddVertex("n5", faasflow.Serializer(serializer))
	maindag.AddModifier("n5", debug("n5"))

	maindag.AddModifier("n6", debug("n6"))

	maindag.AddVertex("n7", faasflow.Serializer(serializer))
	maindag.AddModifier("n7", debug("n7"))

	maindag.AddModifier("n8", debug("n8"))
	maindag.AddModifier("n9", debug("n9"))
	maindag.AddModifier("n10", debug("n10"))

	maindag.AddVertex("n11", faasflow.Serializer(serializer))
	maindag.AddModifier("n11", debug("n11"))

	subdag := faasflow.CreateDag()
	subdag.AddModifier("n13", debug("n13"))
	subdag.AddModifier("n14", debug("n14"))

	superSubDag := faasflow.CreateDag()
	superSubDag.AddModifier("n15.1", debug("n15.1"))
	superSubDag.AddModifier("n15.2", debug("n15.2"))
	superSubDag.AddModifier("n15.3", debug("n15.3"))
	superSubDag.AddModifier("n15.4", debug("n15.4"))
	superSubDag.AddVertex("n15.5", faasflow.Serializer(serializer))
	superSubDag.AddModifier("n15.5", debug("n15.5"))
	superSubDag.AddEdge("n15.1", "n15.2")
	superSubDag.AddEdge("n15.2", "n15.3")
	superSubDag.AddEdge("n15.3", "n15.5")
	superSubDag.AddEdge("n15.1", "n15.4")
	superSubDag.AddEdge("n15.4", "n15.5")

	subdag.AddSubDag("n15", superSubDag)
	subdag.AddVertex("n16", faasflow.Serializer(serializer))
	subdag.AddModifier("n16", debug("n16"))
	subdag.AddEdge("n13", "n14")
	subdag.AddEdge("n13", "n15")
	subdag.AddEdge("n14", "n16")
	subdag.AddEdge("n15", "n16")
	maindag.AddSubDag("n12", subdag)

	maindag.AddVertex("n17", faasflow.Serializer(serializer))
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

	err = flow.ExecuteDag(maindag)

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

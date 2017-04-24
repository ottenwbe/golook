package integration

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/ottenwbe/golook/broker/communication"
	"github.com/ottenwbe/golook/broker/models"
	"github.com/ottenwbe/golook/broker/routing"
	"github.com/ottenwbe/golook/broker/runtime"
	"github.com/ottenwbe/golook/broker/service"
	log "github.com/sirupsen/logrus"
	//"testing"
	"time"
)

/*func TestIntegration(t *testing.T) {
	RunInDockerContainer(runTest)
}*/

func runTest(client *docker.Client, container *docker.Container) {
	var testRouter = routing.NewRouter("myr")

	containerInfo, _ := client.InspectContainer(container.ID)

	log.Infof("connecting rpc client to %s", containerInfo.NetworkSettings.IPAddress)

	testRpcClient := communication.NewLookupRPCClient("http://" + containerInfo.NetworkSettings.IPAddress + ":8080")
	testRouter.NewNeighbor(routing.NewKey(containerInfo.NetworkSettings.MacAddress), testRpcClient)

	st := &models.SystemFiles{System: runtime.GolookSystem, Files: nil}
	res := testRouter.BroadCast(service.SYSTEM_REPORT, service.SystemTransfer{Uuid: runtime.GolookSystem.UUID, System: st, IsDeletion: false})

	log.Println(res)
}

func RunInDockerContainer(f func(client *docker.Client, container *docker.Container)) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Cannot connect to Docker daemon: %s", err)
	}
	container, err := client.CreateContainer(createOptions("golook:server"))
	if err != nil {
		log.Fatalf("Cannot create Docker container; make sure docker daemon is started: %s", err)
	}
	defer func() { //ensure container stops again
		if err := client.RemoveContainer(docker.RemoveContainerOptions{
			ID:    container.ID,
			Force: true,
		}); err != nil {
			log.Fatalf("cannot remove container: %s", err)
		}
	}()

	err = client.StartContainer(container.ID, &docker.HostConfig{})
	if err != nil {
		log.Fatalf("Cannot start Docker container: %s", err)
	}

	// some time for the container
	time.Sleep(time.Second)

	f(client, container)
}

func createOptions(dbname string) docker.CreateContainerOptions {
	ports := make(map[docker.Port]struct{})
	ports["8383"] = struct{}{}
	opts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        dbname,
			ExposedPorts: ports,
		},
	}

	return opts
}

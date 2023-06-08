package nomad

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
)

type Nomad struct {
	client *api.Client
}

func NewNomad() (*Nomad, error) {
	client, err := api.NewClient(&api.Config{})
	if err != nil {
		return nil, err
	}
	return &Nomad{client: client}, nil
}

// pointerOf returns a pointer to a.
func pointerOf[A any](a A) *A {
	return &a
}

// CreatePage starts a new Nomad job, which starts a nginx container
// serving the static page downloaded from the URL. If script is true, the
// script at the URL is executed to generate the page.
// An error is returned if the job could not be started.
// The URL of the page is returned if the job was started successfully.
func (nc *Nomad) CreatePage(name string, URL string, script bool) (string, error) {
	job := nc.createPageJob(name, URL, script)

	_, _, err := nc.client.Jobs().Register(job, nil)
	if err != nil {
		return "", err
	}

	serviceInfo, _, err := nc.client.Services().Get(name, nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%v:%v", serviceInfo[0].Address, serviceInfo[0].Port), nil
}

var scriptEntrypoint = `#!/bin/sh

/bin/chmod +x /generate.sh
/generate.sh > /usr/share/nginx/html/index.html
exec "$@"
`

func (nomad *Nomad) createPageJob(name string, URL string, script bool) *api.Job {
	// Create the job
	job := api.NewServiceJob(name, name, "", 0)

	// Create the task group, add it to the job
	taskGroup := api.NewTaskGroup("main", 1)
	job.TaskGroups = append(job.TaskGroups, taskGroup)

	// Create the service
	service := api.Service{
		Name:      name,
		Provider:  "nomad",
		PortLabel: "http",
	}
	taskGroup.Services = append(taskGroup.Services, &service)

	// Create the task using the docker driver, add it to the task group
	task := api.NewTask("nginx", "docker")
	taskGroup.Tasks = append(taskGroup.Tasks, task)
	// Configure task group network
	port := api.Port{
		Label: "http",
		To:    80,
	}
	taskGroup.Networks = append(taskGroup.Networks, &api.NetworkResource{
		DynamicPorts: []api.Port{port},
	})

	// Setup the task configuration
	task.Config = make(map[string]interface{})
	task.Config["image"] = "nginx"
	task.Config["ports"] = []string{"http"}

	// If script is set, override the docker entrypoint to run the script which
	// generates the page. Otherwise, mount the page as a volume.
	if script {
		task.Templates = append(task.Templates, &api.Template{
			EmbeddedTmpl: pointerOf(scriptEntrypoint),
			DestPath:     pointerOf("local/script-entrypoint.sh"),
			Perms:        pointerOf("555"),
		})

		task.Artifacts = append(task.Artifacts, &api.TaskArtifact{
			GetterSource: &URL,
			GetterMode:   pointerOf("file"),
			RelativeDest: pointerOf("local/generate.sh"),
		})

		task.Config["volumes"] = []string{
			"local/script-entrypoint.sh:/script-entrypoint.sh",
			"local/generate.sh:/generate.sh",
		}

		task.Config["entrypoint"] = []string{"/script-entrypoint.sh"}
		// Reuse the default nginx command and args, which need to be set
		// because the entrypoint is overridden.
		task.Config["command"] = "nginx"
		task.Config["args"] = []string{"-g", "daemon off;"}
	} else {
		task.Artifacts = append(task.Artifacts, &api.TaskArtifact{
			GetterSource: &URL,
			GetterMode:   pointerOf("file"),
			RelativeDest: pointerOf("local/index.html"),
		})
		task.Config["volumes"] = []string{"local/index.html:/usr/share/nginx/html/index.html"}
	}
	return job
}

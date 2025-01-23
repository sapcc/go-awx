# AWX Go Client

This is a Go client for interacting with the AWX API. AWX is an open-source project that provides a web-based user interface, REST API, and task engine built on top of Ansible.

## Installation

To install the AWX Go client, use `go get`:

```sh
go get github.com/sapcc/go-awx
```

## Usage

Here is a basic example of how to use the AWX Go client:

```go
package main

import (
    "fmt"
    "log"
    "github.com/sapcc/go-awx/awx"
)

func main() {
    client, err := awx.NewClient(awx.Options{
        Endpoint: "https://your-awx.cloud/api/v2",
        Username: "username",
		Password: "password",
		Token:    "token",
    })
    if err != nil {
        log.Fatalf("Failed to create AWX client: %v", err)
    }
    input := ListJobsInput{ID: "1"}
	jobs := JobList{}
    err = client.List(context.Background(), awx.ObjectKey{Resource: "jobs"}, &jobs, nil)
    if err != nil {
        log.Fatalf("Failed to list inventories: %v", err)
    }

    for _, job := range jobs.Results {
        fmt.Printf("job: %s\n", job.Name)
    }
}
```

## Features

- List inventories
- Create, update, and delete inventories
- List job templates
- Launch job templates
- And more...

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [AWX Project](https://github.com/ansible/awx)
- [Ansible](https://www.ansible.com/)

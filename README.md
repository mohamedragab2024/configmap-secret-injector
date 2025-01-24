# Golang Kubernetes Operator

This project is a simple Kubernetes operator written in Go that watches for changes in ConfigMaps. It specifically looks for ConfigMaps that contain the annotation `secret-injector/enabled: true` and updates their data accordingly.

## Project Structure

```
golang-k8s-operator
├── bin
│   └── configmap-secret-injector # Compiled binary
├── cmd
│   └── main.go                   # Entry point of the operator
├── internal
│   └── controllers
│       └── configmap_controller.go # Implementation of the ConfigMap controller
├── Dockerfile                    # Instructions for building the Docker image
├── go.mod                        # Module dependencies
├── go.sum                        # Checksums for module dependencies
├── Makefile                      # Commands for building and managing the operator
└── README.md                     # Project documentation
```

## Getting Started

To get started with this operator, follow these steps:

1. **Clone the repository**:
   ```
   git clone https://github.com/mohamedragab2024/configmap-secret-injector.git
   cd golang-k8s-operator
   ```

2. **Build the operator**:
   ```
   make build
   ```

3. **Run the operator**:
   ```
   make run
   ```

## Usage

This operator will watch for ConfigMaps in the Kubernetes cluster. Ensure that your ConfigMaps have the annotation `secret-injector/enabled: true` and specify the secret name with `secret-injector/secret-name` for the operator to take action.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

# SwitchBot Middleware

SwitchBot Middleware is a middleware service that interacts with SwitchBot devices. It provides an API to control SwitchBot devices such as turning them on or off, toggling their state, and more.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Endpoints](#endpoints)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install and run the SwitchBot Middleware, follow these steps:

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/switchbot-middleware.git
    cd switchbot-middleware
    ```

2. Build the Docker image and push it to a container registry:
    ```sh
    docker build -t <registry>/switchbot-middleware .
    docker push <registry>switchbot-middleware
    ```

3. Update the `manifests/deploy.yaml` file with the container image URL:
    ```yaml
    spec:
      containers:
        - name: switchbot-middleware
          image: <registry>/switchbot-middleware
    ```

4. Deploy the application using Kubernetes:

    ```sh
    kubectl apply -f manifests/deploy.yaml
    ```

    or using Helm:

    edit the `manifests/helm/values.yaml` file with the container image URL:

    ```yaml
    image:
      repository: <registry>/switchbot-middleware
    ```

    ```sh
    $ helm install switchbot-middleware manifests/helm
    ```

## Usage

To use the SwitchBot Middleware, you need to have a running instance of the service. You can interact with the API using tools like `curl`.

Example:
```sh
curl http://<your-service-url>/light/toggle
```

## Configuration

The SwitchBot Middleware can be configured using a YAML file. The configuration file should be mounted to the `/app/switchbot-configuration.yaml` file inside the container or specified file path using the `-f` flag when running the application.

YAML configuration structure:
```yaml
switchbot-configuration:
  - name: <device name>
    path: "path/to/request/endpoint"
    type: "turnOn" | "turnOff" | "toggle"
    deviceIds:
      - <device id>
```

Example configuration file:
```yaml
switchbot-configuration:
  - name: "light-turnon"
    path: "/light/turnon"
    type: "turnOn"
    deviceIds:
      - "01-202304012328-87495896"
  - name: "light-turnoff"
    path: "/light/turnoff"
    type: "turnOff"
    deviceIds:
      - "01-202304012328-87495896"
```

for describe the configuration file, see example/config.yaml

## For example

use with streamdeck plugin ninja api plugin like below video

[![IMAGE ALT TEXT HERE](https://img.youtube.com/vi/T_SnHQxt-YI/0.jpg)](https://www.youtube.com/watch?v=T_SnHQxt-YI)
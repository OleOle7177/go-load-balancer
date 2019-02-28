#### Simple Load Balancer

<b>Usage:</b>

1) go build -o lb  
2) ./lb

<b>Config:</b>
place configs in config.yaml

<b>Test:</b>
```curl -XGET 127.0.0.1:5000 -H "Host: header-from-yaml"```

<b>Further work:</b>
 - use heap to prioritize proxy backends
 - tests for config parse
 - read config path from flag
 - tests for proxy requests
 - log requests
 - log errors

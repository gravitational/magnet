# Hello World
`go run mage.go env`

Load env defaults from make.

```
❯ go run mage.go env
arch:  amd64
goVersion:  1.9.2
k8sVersion:  v1.6.6
 ```

Override an env variable as normal
```
❯ ARCH=test go run mage.go env
arch:  test
goVersion:  1.9.2
k8sVersion:  v1.6.6
```

Use the common help target:
```
❯ go run mage.go help:envs
        ENV       | VALUE |   DEFAULT   |                            SHORT DESCRIPTION                              
------------------+-------+-------------+---------------------------------------------------------------------------
  ARCH            |       | amd64       | Set the arch (Default from make)                                          
  DEBIAN_FRONTEND |       |             | Set to noninteractive or stderr to null to enable non-interactive output  
  GO_VERSION      |       | 1.9.2       | Set the golang version (Default from make)                                
  K8S_VERSION     |       | v1.6.6      | Set the k8s version (Default from make)                                   
  XDG_CACHE_HOME  |       | build/cache | Location to store/cache build assets                         
  ```
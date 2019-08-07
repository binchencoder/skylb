# 北京环境

## 北京使用者

### etcd 镜像

```bash
docker pull harbor.eff.com/bld_tools/etcd
```

### skylb 镜像

```bash
docker pull harbor.eff.com/skylb/skylb
```

### skylb web镜像

```bash
docker pull harbor.eff.com/skylb/webserver
```

镜像下载后，在 workspace 下， 执行

```bash
docker-compose -f docker-compose/dev/skylb/docker-compose.yml up
```

## 北京开发者

### etcd 镜像

```bash
bazel run bld_tools/docker/etcd:latest
docker push harbor.eff.com/bld_tools/etcd
```

### skylb 镜像

```bash
bazel run skylb/cmd/skylb:latest.latest
docker push harbor.eff.com/skylb/skylb
```

### skylb web镜像

```bash
bazel run skylb/cmd/webserver:latest
docker push harbor.eff.com/skylb/webserver
```

#!/usr/bin/env bash
# 禁止部署配置文件
rm -rf configs/douyacun.yml
rsync -azv ./* douyacun:${deploy_dir}
ssh douyacun "export _DOUYACUN_CONF='${deploy_dir}/configs/douyacun.yml' && ${deploy_dir}/bin/douyacun stop && ${deploy_dir}/bin/douyacun start"
curl http://localhost:9003/ping

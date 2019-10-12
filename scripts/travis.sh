#!/usr/bin/env bash
# 禁止部署配置文件
rm -rf configs/douyacun.yml
cd .. && rsync -azv api.douyacun.com/ douyacun:${deploy_dir}
ssh douyacun "export ${deploy_dir}/internal/config/douyacun.yml && ${deploy_dir}/bin/douyacun stop && ${deploy_dir}/bin/douyacun start"

# 部署

## 简单部署
```
nohup ./workflow-api start -f /etc/workflow/workflow.toml &> workflow-api.log &
nohup ./workflow-scheduler start -f /etc/workflow/workflow.toml &> workflow-scheduler.log &
nohup ./workflow-node start -f /etc/workflow/workflow.toml &> workflow-node.log &
```
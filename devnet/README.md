# Memory allocation
| container          | measured VmHWM | set limit | service/job  |
|--------------------|----------------|-----------|--------------|
| guardiand          | 2777Mi         | 4000Mi    | service      |
| ganache            | 471Mi          | 500Mi     | service      |
| tests (eth devnet) | ?              | 1Gi       | could be job |
| mine               | 271Mi          | 300Mi     | service      |
| spy                | 116Mi          | 150Mi     | service      |
| algorand-postgres  | ?              | 80Mi      | service      |
| algorand-algod     | 644Mi          | 1000Mi    | service      |
| algorand-indexer   | ?              | 100Mi     | service      |
| algorand-contracts | ?              | 200Mi     | could be job |
| aptos-node         | 143Mi          | 500Mi     | service      |
| aptos-contracts    | ?              | 300Mi     | could be job |
| btc-node           | 310Mi          | 350Mi     | service      |
| near-node          | 639Mi          | 700Mi     | service      |
| near-deploy        | 462Mi          | 500Mi     | could be job |
| solana-devnet      | 1769Mi         | 2000Mi    | service      |
| solana-setup       | ?              | 750Mi     | could be job |
| spy-listener       | ?              | 150Mi     | service      |
| spy-relayer        | 76Mi           | 100Mi     | service      |
| spy-wallet-monitor | 81Mi           | 100Mi     | service      |
| spy                | 102Mi          | 120Mi     | service      |
| terra-terrad       | 343Mi          | 400Mi     | service      |
| terra-contracts    | ?              | 200Mi     | could be job |
| fcd-postgres       | ?              | 50Mi      | service      |
| fcd-collector      | ?              | 500Mi     | service      |
| fcd-api            | ?              | 200Mi     | service      |
| wormchaind         | 559Mi          | 1000Mi    | service      |
| sdk-ci-tests       | ?              | 2500Mi    | job          |
| spydk-ci-tests     | ?              | 1000Mi    | job          |

## Debugging
* oomkill messages should be in `/var/log/messages`.
* kubectl top:
    * Enable metrics server: `minikube addons enable metrics-server`
    * Wait for ~60s (metrics are usually collected every minute)
    * `kubectl top pod --containers=true`
* Get the max memory consumption of a process over its lifetime:
    * e.g. for `algod`: `name=[a]lgod;P=$(ps aux | grep $name | awk '{print $2}'); grep ^VmHWM /proc/$P/status | awk '{print $2}'`
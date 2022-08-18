interface WorkerAction {
  chainId: ChainId;
  id: ActionId;
  data: Object;
  description?: string,
  depedencies?: ActionId[]
}

type ActionId = "UUID"; // todo: import real type

type ContractFilter = {
  emitterAddress: string;
  chainId: ChainId;
};

type EVMToolbox = {};

type SolanaToolbox = {};

type CosmToolbox = {};

//TODO add loggers

//TODO scheduler w/ staging area for when multiple VAAs are rolling in
interface Executor {
  relayEvmAction?: (
    walletToolbox: EVMToolbox,
    action: WorkerAction,
    queuedActions: WorkerAction
  ) => ActionQueueUpdate;
  relaySolanaAction?: (
    walletToolbox: SolanaToolbox,
    action: WorkerAction,
    queuedActions: WorkerAction
  ) => ActionQueueUpdate;
  relayCosmAction?: (
    walletToolbox: CosmToolbox,
    action: WorkerAction,
    queuedActions: WorkerAction
  ) => ActionQueueUpdate;
}

type ActionQueueUpdate = {
  enqueueActions: WorkerAction[];
  removeActionIds: string[];
};


// todo: import from sdk
type ChainId = number
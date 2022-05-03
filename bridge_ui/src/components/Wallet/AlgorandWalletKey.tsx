import { useCallback, useState } from "react";
import { useAlgorandContext } from "../../contexts/AlgorandWalletContext";
import AlgorandConnectWalletDialog from "./AlgorandConnectWalletDialog";
import ToggleConnectedButton from "./ToggleConnectedButton";

const AlgorandWalletKey = () => {
  const { disconnect, accounts } = useAlgorandContext();

  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const connect = useCallback(() => {
    setIsDialogOpen(true);
  }, [setIsDialogOpen]);

  const closeDialog = useCallback(() => {
    setIsDialogOpen(false);
  }, [setIsDialogOpen]);

  return (
    <>
      <ToggleConnectedButton
        connect={connect}
        disconnect={disconnect}
        connected={!!accounts[0]}
        pk={accounts[0]?.address || ""}
      />
      <AlgorandConnectWalletDialog
        isOpen={isDialogOpen}
        onClose={closeDialog}
      />
    </>
  );
};

export default AlgorandWalletKey;

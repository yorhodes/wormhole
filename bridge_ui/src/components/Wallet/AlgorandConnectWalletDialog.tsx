import {
  Dialog,
  DialogTitle,
  Divider,
  IconButton,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  makeStyles,
} from "@material-ui/core";
import CloseIcon from "@material-ui/icons/Close";
import { ConnectType, useWallet } from "@terra-money/wallet-provider";
import WalletConnect from "@walletconnect/client";
import QRCodeModal from "algorand-walletconnect-qrcode-modal";
import { useCallback } from "react";

const useStyles = makeStyles((theme) => ({
  flexTitle: {
    display: "flex",
    alignItems: "center",
    "& > div": {
      flexGrow: 1,
      marginRight: theme.spacing(4),
    },
    "& > button": {
      marginRight: theme.spacing(-1),
    },
  },
  icon: {
    height: 24,
    width: 24,
  },
}));

function connectWalletConnect() {
  // Create a connector
  const connector = new WalletConnect({
    bridge: "https://bridge.walletconnect.org", // Required
    qrcodeModal: QRCodeModal,
  });

  // Check if connection is already established
  if (!connector.connected) {
    console.log("connected");
    // create new session
    connector.createSession();
  }

  // Subscribe to connection events
  connector.on("connect", (error, payload) => {
    console.log("connect");
    if (error) {
      throw error;
    }

    // Get provided accounts
    const { accounts } = payload.params[0];
  });

  connector.on("session_update", (error, payload) => {
    console.log("update");
    if (error) {
      throw error;
    }

    // Get updated accounts
    const { accounts } = payload.params[0];
    console.log(accounts);
  });

  connector.on("disconnect", (error, payload) => {
    console.log("disconnect");
    if (error) {
      throw error;
    }
  });
}

const WalletOptions = ({
  type,
  identifier,
  connect,
  onClose,
  icon,
  name,
}: {
  type: ConnectType;
  identifier: string;
  connect: (
    type: ConnectType | undefined,
    identifier: string | undefined
  ) => void;
  onClose: () => void;
  icon: string;
  name: string;
}) => {
  const classes = useStyles();

  const handleClick = useCallback(() => {
    connect(type, identifier);
    onClose();
  }, [connect, onClose, type, identifier]);
  return (
    <ListItem button onClick={handleClick}>
      <ListItemIcon>
        <img src={icon} alt={name} className={classes.icon} />
      </ListItemIcon>
      <ListItemText>{name}</ListItemText>
    </ListItem>
  );
};

const AlgorandConnectWalletDialog = ({
  isOpen,
  onClose,
}: {
  isOpen: boolean;
  onClose: () => void;
}) => {
  const classes = useStyles();
  return (
    <Dialog open={isOpen} onClose={onClose}>
      <DialogTitle>
        <div className={classes.flexTitle}>
          <div>Select your wallet</div>
          <IconButton onClick={onClose}>
            <CloseIcon />
          </IconButton>
        </div>
      </DialogTitle>
      <List>
        <ListItem button onClick={connectWalletConnect}>
          {/* <ListItemIcon>
            <img src={icon} alt={name} className={classes.icon} />
          </ListItemIcon> */}
          <ListItemText>WalletConnect</ListItemText>
        </ListItem>
      </List>
    </Dialog>
  );
};

export default AlgorandConnectWalletDialog;

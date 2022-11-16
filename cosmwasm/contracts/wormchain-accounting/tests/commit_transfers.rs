mod helpers;

use accounting::state::{transfer, Transfer};
use cosmwasm_std::{to_binary, Uint256};
use helpers::*;
use wormhole_bindings::Signature;

pub fn set_up(count: usize) -> Vec<Transfer> {
    let mut out = Vec::with_capacity(count);
    for i in 0..count {
        let key = transfer::Key::new(i as u16, [i as u8; 32].into(), i as u64);
        let data = transfer::Data {
            amount: Uint256::from(i as u128),
            token_chain: i as u16,
            token_address: [(i + 1) as u8; 32].into(),
            recipient_chain: (i + 2) as u16,
        };

        out.push(Transfer { key, data });
    }

    out
}

#[test]
fn basic() {
    const COUNT: usize = 7;

    let (wh, mut contract) = proper_instantiate(Vec::new(), Vec::new(), Vec::new());
    let index = wh.guardian_set_index();

    let txs = set_up(COUNT);
    let transfers = to_binary(&txs).unwrap();
    let signatures = wh.sign(&transfers);

    contract
        .commit_transfers(transfers, index, signatures)
        .unwrap();

    for t in txs {
        assert_eq!(t.data, contract.query_transfer(t.key).unwrap());
    }
}

#[test]
fn invalid_transfer() {
    const COUNT: usize = 7;

    let (wh, mut contract) = proper_instantiate(Vec::new(), Vec::new(), Vec::new());
    let index = wh.guardian_set_index();

    let mut txs = set_up(COUNT);
    txs[3].data.token_chain = txs[3].data.recipient_chain;

    let transfers = to_binary(&txs).unwrap();
    let signatures = wh.sign(&transfers);

    contract
        .commit_transfers(transfers, index, signatures)
        .expect_err("successfully committed transfer list containing invalid transfer");
}

#[test]
fn no_quorum() {
    const COUNT: usize = 7;

    let (wh, mut contract) = proper_instantiate(Vec::new(), Vec::new(), Vec::new());
    let index = wh.guardian_set_index();
    let quorum = wh
        .calculate_quorum(index, contract.app().block_info().height)
        .unwrap() as usize;

    let txs = set_up(COUNT);
    let transfers = to_binary(&txs).unwrap();
    let mut signatures = wh.sign(&transfers);
    signatures.truncate(quorum - 1);

    contract
        .commit_transfers(transfers, index, signatures)
        .expect_err("successfully committed transfer list without a quorum of signatures");
}

#[test]
fn bad_serialization() {
    const COUNT: usize = 7;

    let (wh, mut contract) = proper_instantiate(Vec::new(), Vec::new(), Vec::new());
    let index = wh.guardian_set_index();

    let txs = set_up(COUNT);

    // Rather than serializing a Vec, just serialize a single element.
    let transfers = to_binary(&txs[0]).unwrap();
    let signatures = wh.sign(&transfers);

    contract
        .commit_transfers(transfers, index, signatures)
        .expect_err("successfully committed transfer with bad serialization");
}

#[test]
fn bad_signature() {
    const COUNT: usize = 7;

    let (wh, mut contract) = proper_instantiate(Vec::new(), Vec::new(), Vec::new());
    let guardian_set_index = wh.guardian_set_index();
    let quorum = wh
        .calculate_quorum(guardian_set_index, contract.app().block_info().height)
        .unwrap() as usize;

    let txs = set_up(COUNT);
    let transfers = to_binary(&txs).unwrap();
    let mut signatures = wh.sign(&transfers);
    signatures.truncate(quorum - 1);

    // Flip a bit in the signature so it becomes invalid.
    let Signature { index, signature } = signatures.swap_remove(0);
    let mut v = Vec::from(signature);
    v[0] ^= 1;
    signatures.push(Signature {
        index,
        signature: v.into(),
    });

    contract
        .commit_transfers(transfers, guardian_set_index, signatures)
        .expect_err("successfully committed transfer with bad signature");
}

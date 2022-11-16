module token_bridge::wrapped {
    use sui::tx_context::TxContext;
    //use sui::object::{Self, UID};
    use sui::coin::{Self};
    //use sui::coin::{Self, Coin, TreasuryCap};
    //use sui::transfer::{Self};

    use token_bridge::bridge_state::{Self, BridgeState};
    use token_bridge::vaa::{Self as token_bridge_vaa};
    use token_bridge::asset_meta::{AssetMeta, Self, get_decimals};
    use token_bridge::treasury::{Self};

    use wormhole::state::{Self as state, State as WormholeState};
    use wormhole::myvaa::{Self as corevaa};

    const E_WRAPPING_NATIVE_COIN: u64 = 0;
    const E_WRAPPING_REGISTERED_NATIVE_COIN: u64 = 1;
    const E_WRAPPED_COIN_ALREADY_INITIALIZED: u64 = 2;

    public entry fun create_wrapped_coin<CoinType: drop>(
        state: &mut WormholeState,
        bridge_state: &mut BridgeState,
        vaa: vector<u8>,
        witness: CoinType,
        ctx: &mut TxContext,
    ) {
        let vaa = token_bridge_vaa::parse_verify_and_replay_protect(state, bridge_state, vaa, ctx);
        let asset_meta: AssetMeta = asset_meta::parse(corevaa::destroy(vaa));
        let decimals = get_decimals(&asset_meta);
        let treasury_cap = coin::create_currency<CoinType>(witness, decimals, ctx);
        // assert emitter is registered

        // TODO (pending Mysten Labs uniform token standard) -  extract/store decimals, token name, symbol, etc. from asset meta

        let t_cap_store = treasury::create_treasury_cap_store<CoinType>(treasury_cap, ctx);

        let origin_chain = asset_meta::get_token_chain(&asset_meta);
        assert!(origin_chain != state::get_chain_id(state), E_WRAPPING_NATIVE_COIN);
        assert!(!bridge_state::is_registered_native_asset<CoinType>(bridge_state), E_WRAPPING_REGISTERED_NATIVE_COIN);
        assert!(!bridge_state::is_wrapped_asset<CoinType>(bridge_state), E_WRAPPED_COIN_ALREADY_INITIALIZED);

        let external_address = asset_meta::get_token_address(&asset_meta);
        let wrapped_asset_info = bridge_state::create_wrapped_asset_info(origin_chain, external_address, ctx);

        bridge_state::register_wrapped_asset<CoinType>(bridge_state, wrapped_asset_info);
        bridge_state::store_treasury_cap<CoinType>(bridge_state, t_cap_store);
    }
}
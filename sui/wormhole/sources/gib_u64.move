module wormhole::gib_u64{
    public fun gib_u64(_: u32): u64{
        let x: u64 = 3;
        return x
    }
}

#[test_only]
module wormhole::test_gib_u64{
    use wormhole::gib_u64::gib_u64;

    #[test]
    public fun test_gib_u64(){
        let _z = (123 as u16);
        let _w = (_z as u32);
        let _y = (_w as u64);
        let _x = gib_u64(3);

        let res: u128 = 0;
        let i = 0;
        while (i < 16) {
            res = (res << 8) + (3 as u128);
            i = i + 1;
        };
    }
}


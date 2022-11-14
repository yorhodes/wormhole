module wormhole::gib_u32{
    public fun gib_u32(): u32{
        let x: u32 = 3;
        return x
    }
}

#[test_only]
module wormhole::test_gib_u32{

}

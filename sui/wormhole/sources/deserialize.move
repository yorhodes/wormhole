module wormhole::deserialize {
    use std::vector::{Self};
    use wormhole::cursor::{Self, Cursor};

    public fun deserialize_u8(cur: &mut Cursor<u8>): u8 {
        cursor::poke(cur)
    }

    public fun deserialize_u16(cur: &mut Cursor<u8>): u16 {
        let res: u16 = 0;
        let i = 0;
        while (i < 2) {
            let b = cursor::poke(cur);
            res = (res << 8) + (b as u16);
            i = i + 1;
        };
        res
    }

    public fun deserialize_u32(cur: &mut Cursor<u8>): u32 {
        let res: u32 = 0;
        let i = 0;
        while (i < 4) {
            let b = cursor::poke(cur);
            res = (res << 8) + (b as u32);
            i = i + 1;
        };
        res
    }

    public fun deserialize_u64(cur: &mut Cursor<u8>): u64 {
        let res: u64 = 0;
        let i = 0;
        while (i < 8) {
            let b = cursor::poke(cur);
            res = (res << 8) + (b as u64);
            i = i + 1;
        };
        res
    }

    public fun deserialize_u128(cur: &mut Cursor<u8>): u128 {
        let res: u128 = 0;
        let i = 0;
        while (i < 16) {
            let b = cursor::poke(cur);
            res = (res << 8) + (b as u128);
            i = i + 1;
        };
        res
    }

    public fun deserialize_u256(cur: &mut Cursor<u8>): u256 {
        let res: u256 = 0;
        let i = 0;
        while (i < 32) {
            let b = cursor::poke(cur);
            res = (res << 8) + (b as u256);
            i = i + 1;
        };
        res
    }

    public fun deserialize_vector(cur: &mut Cursor<u8>, len: u64): vector<u8> {
        let result = vector::empty();
        while (len > 0) {
            vector::push_back(&mut result, cursor::poke(cur));
            len = len - 1;
        };
        result
    }

}

#[test_only]
module wormhole::deserialize_test {
    use wormhole::cursor;
    use wormhole::deserialize::{
        deserialize_u8,
        deserialize_u16,
        deserialize_u32,
        deserialize_u64,
        deserialize_u128,
        deserialize_vector,
    };

    #[test]
    fun test_deserialize_u8() {
        let cur = cursor::cursor_init(x"99");
        let byte = deserialize_u8(&mut cur);
        assert!(byte==0x99, 0);
        cursor::destroy_empty(cur);
    }

    #[test]
    fun test_deserialize_u16() {
        let cur = cursor::cursor_init(x"9987");
        let u = deserialize_u16(&mut cur);
        assert!(u == 0x9987, 0);
        cursor::destroy_empty(cur);
    }

    #[test]
    fun test_deserialize_u32() {
        let cur = cursor::cursor_init(x"99876543");
        let u = deserialize_u32(&mut cur);
        assert!(u == 0x99876543, 0);
        cursor::destroy_empty(cur);
    }

    #[test]
    fun test_deserialize_u64() {
        let cur = cursor::cursor_init(x"1300000025000001");
        let u = deserialize_u64(&mut cur);
        assert!(u == 0x1300000025000001, 0);
        cursor::destroy_empty(cur);
    }

    #[test]
    fun test_deserialize_u128() {
        let cur = cursor::cursor_init(x"130209AB2500FA0113CD00AE25000001");
        let u = deserialize_u128(&mut cur);
        assert!(u == 0x130209AB2500FA0113CD00AE25000001, 0);
        cursor::destroy_empty(cur);
    }

    #[test]
    fun test_deserialize_vector() {
        let cur = cursor::cursor_init(b"hello world");
        let hello = deserialize_vector(&mut cur, 5);
        deserialize_u8(&mut cur);
        let world = deserialize_vector(&mut cur, 5);
        assert!(hello == b"hello", 0);
        assert!(world == b"world", 0);
        cursor::destroy_empty(cur);
    }

}

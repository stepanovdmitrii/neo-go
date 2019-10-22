# Changelog

This document outlines major changes between releases.

## 0.51.0 "Confirmation" (17 Oct 2019)

With over a 100 commits made since 0.50.0 release 0.51.0 brings with it full
block verification, improved and fixed transaction verifications,
implementation of most of interop functions and VM improvements. Block
verification is an important milestone on the road to full consensus node
support and it required implementing a lot of other associated functionality.

New features:
 * CHECKSIG, VERIFY and CHECKMULTISIG instructions in VM (#269)
 * witness verification logic for transactions (#368)
 * NEWMAP, HASKEY, KEYS and VALUES instructions, support for Map type in
   PICKITEM, SETITEM, REMOVE, EQUAL, ARRAYSIZE (#359)
 * configurable transaction verification on block addition (#415, #418)
 * contract storage and support for VM to call contracts via APPCALL/TAILCALL
   (#417)
 * support for Interop type in VM (#417)
 * VM now has stepInto/stepOver/stepOut method implementations for debugging
   matching neo-vm behavior (#187)
 * storage support for contracts (#418)
 * added around 90% of interop functions (#418)
 * Invocation TX processing now really does invoke contracts using internal VM
   (#418)
 * blocks are now completely verified when added to the chain (if not
   configured otherwise; #12, #418)

Behavior changes:
 * full block verification is now enabled for all network types
 * block's transaction verification enabled for privnet setups, mainnet and
   testnet don't have it enabled

Technical improvements:
 * GetVarIntSize and GetVarStringSize were removed from the io package (use
   GetVarSize instead; #408)
 * OVER implementation was optimized to not pop the top element from the stack
   (#406, part of #196 work)
 * vm.VM was extended with HasFailed() method to check its state (previously
   external VM users couldn't do it; #411)
 * redesigned input block queue mechanism, now it's completely moved out of
   the Blockchain, which only accepts the next block via AddBlock() (#414)
 * unpersisted blocks are now fully available with the Blockchain (thus we
   have symmetry now in AddBlock/GetBlock APIs; #414, #366)
 * removed duplicating batch structures from BoltDB and Redis code, now all of
   them use the same batch as MemoryStore does (#414)
 * MemoryStore was exporting its mutex for no good reason, now it's hidden
   (#414)
 * storage layer now returns ErrKeyNotFound for all DBs in appropriate
   situations (#414)
 * VM's PopResult() now doesn't panic if there is no result (#417)
 * VM's Element now has a Value() method to quickly get the item value (#417)
 * VM's stack PushVal() method now accepts uint8/16/32/64 (#417, #418)
 * VM's Element now has TryBool() method similar to Bool(), but without a
   panic (for external VM users; #417)
 * VM has now completely separated instruction read and execution phases
   (#417)
 * Store interface now has Delete method (#418)
 * Store tests were reimplemented to use one test set for all Store
   implementations, including LevelDB that was not tested at all previously
   (#418)
 * Batch interface doesn't have Len method now as it's not used at all (#418)
 * New*FromRawBytes functions were renamed to New*FromASN1 in the keys
   package, previous naming made it easy to confuse them with functions
   operating with NEO serialization format (#418)
 * PublicKey's IsInfinity method is exported now (#418)
 * smartcontract package now has CreateSignatureRedeemScript() matching C#
   code (#418)
 * vm package now has helper functions
   IsSignatureContract/IsMultiSigContract/IsStandardContract matching C# code
   (#418)
 * Blockchain's GetBlock() now returns full block with transactions (#418)
 * Block's Verify() was changed to return specific error (#418)
 * Blockchain's GetTransationResults was renamed into GetTransactionResults
   (#418)
 * Blockchainer interface was extended with GetUnspentCoinState,
   GetContractState and GetScriptHashesForVerifying methods (#418)
 * introduced generic MemCacheStore that is used now for write caching
   (including temporary stores for transaction processing) and batched
   persistence (#425)

Bugs fixed:
 * useless persistence failure message printed with no error (#409)
 * persistence error message being printed twice (#409)
 * segmentation fault upon receival of message that is not currently handled
   properly (like "consensus" message; #409)
 * BoltDB's Put for a batch wasn't copying data which could lead to data
   corruption (#409)
 * APPEND instruction applied to struct element was not copying it like neo-vm
   does (#405, part of #196 work)
 * EQUAL instruction was comparing array contents, while it should've compared
   references (#405, part of #196 work)
 * SUBSTR instruction was failing for out of bounds length parameters while it
   should've truncated them to string length (#406, part of #196 work)
 * SHL and SHR implementations had no limits, neo-vm restricts them to
   -256/+256 (#406, part of #196 work)
 * minor VM state mismatches with neo-vm on failures (#405, #406)
 * deadlock on Blockchain init when headers pointer is not in sync with the
   hashes list (#414)
 * node failed to request blocks when headers list was exactly one position
   ahead of block count (#414)
 * TestRPC/getassetstate_positive failed occasionally (#410)
 * panic on block verification with no transactions inside (#415)
 * DutyFlag check in GetScriptHashesForVerifying was not done correctly (#415)
 * default asset expiration for assets created with Register TX was wrong, now
   it matches C# code (#415)
 * Claim transactions needed specific GetScriptHashesForVerifying logic to
   be verified correctly (#415)
 * VerifyWitnesses wasn't properly sorting hashes and witnesses (#415)
 * transactions referring to two outputs of some other transaction were
   failing to verify (#415)
 * wrong program dumps (#295)
 * potential data race in logging code (#418)
 * bogus port check during handshake (#432)
 * missing max size checks in NEWARRAY, NEWSTRUCT, APPEND, PACK, SETITEM
   (#427, part of #373)

## 0.50.0 "Consolidation" (19 Sep 2019)

The first release from the new team focuses on bringing all related
development effort into one codebase, refactoring things, fixing some
long-standing bugs and adding new functionality. This release merges two
radically different development branches --- `dev` and `master` that were
present in the project (along with all associated pull requests) and also
brings in changes made to the compiler in the neo-storm project.

New features:
 * configurable storage backends supporting LevelDB, in-memory DB (for
   testing) and Redis
 * BoltDB support for storage backend
 * updated and extended interop APIs (thanks to neo-storm)

Notable behavior changes:
 * the default configuration for privnet was changed to use ports 20331 and
   20332 so that it doesn't clash with the default dockerized neo-privnet
   setups
 * the default configuration path was changed from `../config` to `./config`,
   at this stage it makes life a bit easier for development, later this will
   be changed to some sane default for production version
 * VM CLI now supports type inference for `run` parameters (you don't need to
   specify types explicitly in most of the cases) and treats `operation`
   parameter as mandatory if anything is passed to `run`
 
VM improvements:
 * added implementation for `EQUAL`, `NZ`, `PICK`, `TUCK`, `XDROP`, `INVERT`,
   `CAT`, `SUBSTR`, `LEFT`, `RIGHT`, `UNPACK`, `REVERSE`, `REMOVE`
 * expanded tests
 * better error messages for different erroneous code
 * implemented item conversions following neo-vm behavior: array to/from
   struct, bigint to/from boolean, anything to bytearray and anything to
   boolean
 * improved compatibility with neo-vm (#394)

Technical improvements:
 * switched to Go 1.12+
 * gofmt, golint (#213)
 * fixed and improved CircleCI builds
 * removed internal rfc6969 package (#285)
 * refactored util/crypto/io packages, removed a lot of duplicating code
 * updated READMEs and user-level documents
 * update Makefile with useful targets
 * dropped internal base58 implementation (#355)
 * updated default seed lists for mainnet and testnet from neo-cli

Bugs fixed:
 * a lot of compiler fixes from neo-storm
 * data access race in memory-backed storage backend (#313)
 * wrong comparison opcode emitted by compiler (#294)
 * decoding error in `publish` transactions (#179)
 * decoding error in `invocation` transactions (#173)
 * panic in `state` transaction decoding
 * double VM run from CLI (#96)
 * non-constant time crypto (#245)
 * APPEND pushed value on the stack and worked for bytearrays (#391)
 * reading overlapping hash blocks from the DB leading to blockchain state
   neo-go couldn't recover from (#393)
 * codegen for `append()` wasn't type-aware and emitted wrong code (#395)
 * node wasn't trying to reconnect to other node if connection failed (#390)
 * stricly follow handshare procedure (#396)
 * leaked connections if disconnect happened before handshake completed (#396)

### Inherited unreleased changes

Some changes were also done before transition to the new team, highlights are:
 * improved RPC testing
 * implemented `getaccountstate`, `validateaddress`, `getrawtransaction` and
   `sendrawtransaction` RPC methods in server
 * fixed `getaccountstate` RPC implementation
 * implemented graceful blockchain shutdown with proper DB closing

## 0.45.14 (not really released, 05 Dec 2018)

This one can technically be found in the git history and attributed to commit
fa1da2cb917cf4dfccbe49d44c5741eec0e0bb65, but it has no tag in the repository
and so can't count as a properly released thing. Still it can be marked as a
point in history with the following changes relative to 0.44.10:
 * switched to Go modules for dependency management
 * various optimizations for basic structures like Uin160/Uint256/Fixed8
 * improved testing
 * added support for `invoke` method in RPC client
 * implemented `getassetstate` in RPC server
 * fixed NEP2Encrypt and added tests
 * added `NewPrivateKeyFromRawBytes` function to the `wallet` package

## 0.44.10 (27 Sep 2018)

This is the last one tagged in the repository, so it's considered as the last
one properly released before 0.50+. Releases up to 0.44.10 seem to be made in
automated fashion and you can find their [changes on
GitHub](https://github.com/nspcc-dev/neo-go/releases).
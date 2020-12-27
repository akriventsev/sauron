pragma solidity 0.4.24;

import "https://github.com/smartcontractkit/chainlink/evm-contracts/src/v0.4/ChainlinkClient.sol";
import "https://github.com/smartcontractkit/chainlink/evm-contracts/src/v0.4/vendor/Ownable.sol";

contract ATestnetConsumer is ChainlinkClient, Ownable {
  uint256 constant private ORACLE_PAYMENT = 1 * LINK;

  //uint256 public currentPrice;
  //int256 public changeDay;
  //bytes32 public lastMarket;

  uint256 public stackedUSDNBallance;
  uint256 public ethUSDNBallance;

  event RequestEthereumUSDNBallance(
    bytes32 indexed requestId,
    uint256 indexed ballance
  );

  event RequestStackedUSDNBallance(
    bytes32 indexed requestId,
    uint256 indexed ballance
  );

  
  constructor() public Ownable() {
    setPublicChainlinkToken();
  }

  function requestEthereumUSDNBallance(address _oracle, string _jobId)
    public
  {
    Chainlink.Request memory req = buildChainlinkRequest(stringToBytes32(_jobId), this, this.fulfillEthereumUSDNBallance.selector);
    req.add("get", "https://api.etherscan.io/api?module=stats&action=tokensupply&contractaddress=0x674C6Ad92Fd080e4004b2312b45f796a192D27a0&apikey=B4RE54AMA6CR3I7MX65R12KR8VFKRKIE9S");
    req.add("path", "result");
    sendChainlinkRequestTo(_oracle, req, ORACLE_PAYMENT);
  }

  function fulfillEthereumUSDNBallance(bytes32 _requestId, uint256 _ballance)
    public
    recordChainlinkFulfillment(_requestId)
  {
    emit RequestEthereumUSDNBallance(_requestId, _ballance);
    ethUSDNBallance = _ballance;
  }

  function requestStackedUSDNBallance(address _oracle, string _jobId)
    public
  {
    Chainlink.Request memory req = buildChainlinkRequest(stringToBytes32(_jobId), this, this.fulfillStackedUSDNBallance.selector);
    req.add("get", "https://nodes.wavesnodes.com/addresses/data/3PNikM6yp4NqcSU8guxQtmR5onr2D4e8yTJ/rpd_balance_DG2xFkPdDwKUoBkzGAhQtLpSGzfXLiCYPEzeKH2Ad24p_3P7RhLuvncw74sinqGa7SvZYgejXxs5gVyk");
    req.add("path", "value");
    sendChainlinkRequestTo(_oracle, req, ORACLE_PAYMENT);
  }

  function fulfillStackedUSDNBallance(bytes32 _requestId, uint256 _stackedUSDNBallance)
    public
    recordChainlinkFulfillment(_requestId)
  {
    emit RequestStackedUSDNBallance(_requestId, _stackedUSDNBallance);
    stackedUSDNBallance = _stackedUSDNBallance;
  }
  
  function getChainlinkToken() public view returns (address) {
    return chainlinkTokenAddress();
  }
  
  function withdrawLink() public onlyOwner {
    LinkTokenInterface link = LinkTokenInterface(chainlinkTokenAddress());
    require(link.transfer(msg.sender, link.balanceOf(address(this))), "Unable to transfer");
  }

  function cancelRequest(
    bytes32 _requestId,
    uint256 _payment,
    bytes4 _callbackFunctionId,
    uint256 _expiration
  )
    public
    onlyOwner
  {
    cancelChainlinkRequest(_requestId, _payment, _callbackFunctionId, _expiration);
  }

  function stringToBytes32(string memory source) private pure returns (bytes32 result) {
    bytes memory tempEmptyStringTest = bytes(source);
    if (tempEmptyStringTest.length == 0) {
      return 0x0;
    }

    assembly { // solhint-disable-line no-inline-assembly
      result := mload(add(source, 32))
    }
  }

}
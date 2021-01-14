# Proof-of-Reserves (PoR) system prototype for algorithmic stablecoin Neutrino USD (pegged to $USDN)

For wrapped tokens, a system is required that monitors and matches balances on the SuSy LU port with those of the wrapped tokens. If there is a discrepancy of some kind, an Emergency Call should be initiated which should record an entry in a smart contract on one of the target chains. This repository contains source code of the service and the contracts. Uses ChainLink’s oracle node.

The general idea behind the PoR (Proof-of-Reserves) concept is to prove that the number of assets on one account matches the number of assets on another. The original blockchain network where USDN in issued is Waves but a large number of tokens were issued as ERC20 through a waves.exchange gateway. The release of USDN on the Ethereum network is accompanied by locking an equivalent amount of USDN on a gateway account on the Waves network. Accordingly, this external service (oracle) monitors the balance of the gateway account, which contains $USDNs, and verifies that exactly the same number of $USDN tokens exist on the Ethereum network.

Tracking information in the Ethereum blockchain is implemented as a smart contract on the test network using chainlink oracles.

## Service cycle
* Every 10 seconds registering the difference via Waves API & Etherscan.io 
* Every 6 hours checking the difference between balances and sending notifications
* Registering the difference if the discrepancy between balances exceeds the conditional threshold of 2%.

![Workflow](PoR-workflow.png?raw=true "Title")

import { TezBridgeSigner } from '@taquito/tezbridge-signer';
import { Tezos } from '@taquito/taquito';

const mountProvider = function () {
  Tezos.setProvider({
    rpc: process.env.TESTNET,
    signer: new TezBridgeSigner()
  });

  console.log('Tezos =>', Tezos);
};

const getBalance = async (address) => {
  let worker = await Tezos.tz.getBalance(address);
  return worker;
};

const getContractInstance = async (contract) => {
  let contractInstance = await Tezos.contract.at(contract);
  return contractInstance;
};

const _longToByteArray = function(long) {
  // Represent input as 8-bit array
  let byteArray = [0, 0, 0, 0, 0, 0, 0, 0];
  for (let index = 0; index < byteArray.length; index ++) {
    let byte = long & 0xff;
    byteArray [ index ] = byte;
    long = (long - byte) / 256 ;
  }
  return byteArray;
};

const _bytesToHex = function (byteArray) {
  return Array.prototype.map.call(byteArray, function(byte) {
    return ('0' + (byte & 0xFF).toString(16)).slice(-2);
  }).join('');
}

/**
 * Message signing - authenticates game session login, propose rule, vote
 * @param {String} msg 
 * @see https://tezostaquito.io/typedoc/classes/_taquito_tezbridge_signer.tezbridgesigner.html
 * 
 * TezBridgeSigner.sign(bytes: string, _watermark?: Uint8Array): Promise<
 *  { bytes: string; prefixSig: any; sbytes: string; sig: any }
 * >
 */
const signMessage = async (msg) => {
  if (typeof msg !== 'number')
    return false;
  // Long to bytes
  let bytes = _longToByteArray(msg);
  // Bytes to hex
  let hex = _bytesToHex(bytes);
  // Sign message
  let signedMsg = await Tezos.signer.sign(hex);
  // console.log('Message Signer', t);
  return signedMsg;
};

module.exports = {
  Tezos: Tezos,
  mountProvider: mountProvider,
  getBalance: getBalance,
  getContractInstance: getContractInstance,
  signMessage: signMessage
};
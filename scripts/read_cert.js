// ============================================================================================================================
// 													Instantiate Chaincode
// This file shows how to instantiate chaincode onto a Hyperledger Fabric Channel via the SDK + FC Wrangler
// ============================================================================================================================
const fs = require('fs');
var winston = require('winston');								//logger module
var path = require('path');
var logger = new (winston.Logger)({
	level: 'debug',
	transports: [
		new (winston.transports.Console)({ colorize: true }),
	]
});

// --- Set Details Here --- //
var config_file = 'marbles_local.json';							//set config file name
var chaincode_id = 'blockcertsx';									//use same ID during the INSTALL proposal
var chaincode_ver = 'v1';										//use same version during the INSTALL proposal

//  --- Use (optional) arguments if passed in --- //
var args = process.argv.slice(2);
if (args[0]) {
	config_file = args[0];
	logger.debug('Using argument for config file', config_file);
}
if (args[1]) {
	chaincode_id = args[1];
	logger.debug('Using argument for chaincode id');
}
if (args[2]) {
	chaincode_ver = args[2];
	logger.debug('Using argument for chaincode version');
}

var cp = require(path.join(__dirname, '../utils/connection_profile_lib/index.js'))(config_file, logger);			//set the config file name here
var fcw = require(path.join(__dirname, '../utils/fc_wrangler/index.js'))({ block_delay: cp.getBlockDelay() }, logger);

console.log('---------------------------------------');
logger.info('Lets instantiate some chaincode -', chaincode_id, chaincode_ver);
console.log('---------------------------------------');
logger.warn('Note: the chaincode should have been installed before running this script');

logger.info('First we enroll');
fcw.enrollWithAdminCert(cp.makeEnrollmentOptionsUsingCert(), function (enrollErr, enrollResp) {
	if (enrollErr != null) {
		logger.error('error enrolling', enrollErr, enrollResp);
	} else {

		console.log('---------------------------------------');
		logger.info('Lets try Reading');
		console.log('---------------------------------------');

		var opts3 = {
			chaincode_id: chaincode_id,
			cc_function: "readCert",
			cc_args: ["myCert"]
		};

		fcw.query_chaincode(enrollResp , opts3, function (err, resp) {
			fs.writeFile('../downloads/response.json', JSON.stringify(resp.parsed), function (err, resp){
				console.log('Writing File');


			});
			console.log('---------------------------------------');
			logger.info('Invoke Read Cert done. Errors:', (!err) ? 'nope' : err);
			console.log('---------------------------------------');
		});

}});

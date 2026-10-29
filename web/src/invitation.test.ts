import test from 'node:test';
import assert from 'node:assert/strict';
import {invitationLocation,invitationIsLoopback} from './invitation.ts';
test('invitation credentials are separated from navigation before normal reads',()=>{
 const token='synthetic-one-use-token';
 const result=invitationLocation({hash:'#invite='+token,pathname:'/',search:'?layout=force'});
 assert.deepEqual(result,{token,cleanURL:'/?layout=force'});
 assert.ok(!result!.cleanURL.includes(token));assert.ok(!result!.cleanURL.includes('#'));
});
test('an empty invitation still has its fragment cleared and is handled as unavailable',()=>{
 assert.deepEqual(invitationLocation({hash:'#invite=',pathname:'/review',search:''}),{token:'invalid',cleanURL:'/review'});
 assert.equal(invitationLocation({hash:'#main',pathname:'/',search:''}),null);
});
test('unrelated fragment values never become a redirect or credential',()=>{
 assert.deepEqual(invitationLocation({hash:'#invite=synthetic&next=https://outside.example/',pathname:'/',search:''}),{token:'synthetic',cleanURL:'/'});
 assert.equal(invitationLocation({hash:'#other=synthetic',pathname:'/',search:''}),null);
});
test('local-access guidance is limited to actual loopback hostnames',()=>{
 for(const url of ['http://127.0.0.1:8096/#invite=synthetic','http://localhost:8096/','http://[::1]:8096/'])assert.equal(invitationIsLoopback(url),true);
 for(const url of ['https://sentrame.example/','http://192.168.1.10:8096/','https://127.0.0.1.example/','https://localhost.example/','invalid'])assert.equal(invitationIsLoopback(url),false);
});

import test from 'node:test';
import assert from 'node:assert/strict';
import {receivedMappingContexts,initialReceivedMapping,applyReceivedContext} from './mappingEvidence.ts';
const context=(service:string,source='browser',signal='logs',carrier='log')=>({name:'pwa.app.opened',service,source,signal,carrier,count:2});
test('new mapping uses the sole received service and producer without inventing a condition',()=>{
 const mapping=initialReceivedMapping('pwa.app.opened',[context('wfs-pwa-browser')]);
 assert.equal(mapping.service,'wfs-pwa-browser');assert.equal(mapping.source,'browser');assert.equal(mapping.carrier,'log');assert.equal(mapping.key,'');assert.equal(mapping.value,'');
});
test('multiple received contexts require a choice and never silently select another service or authority',()=>{
 const contexts=[context('wfs-pwa-browser'),context('api','server'),{...context('ignored'),name:'other'},context('metrics','server','metrics','metric')];
 assert.equal(receivedMappingContexts('pwa.app.opened',contexts).length,2);
 const mapping=initialReceivedMapping('pwa.app.opened',contexts);
 assert.equal(mapping.service,'');assert.equal(mapping.source,'');
 const next=applyReceivedContext({...mapping,step:'App opened',key:'stage.name',value:'search'},contexts[1]);
 assert.equal(next.source,'server');assert.equal(next.service,'api');assert.equal(next.step,'App opened');assert.equal(next.key,'');assert.equal(next.value,'');
});
test('projection variants of the same context combine while absent history keeps selectors broad',()=>{
 const contexts=[context('browser'),context('browser')];
 assert.equal(receivedMappingContexts('pwa.app.opened',contexts)[0].count,4);
 assert.equal(initialReceivedMapping('pwa.app.opened',contexts).service,'browser');
 assert.equal(initialReceivedMapping('missing',contexts).service,'');
});

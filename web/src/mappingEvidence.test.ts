import test from 'node:test';
import assert from 'node:assert/strict';
import {receivedMappingContexts,initialReceivedMapping,applyReceivedContext} from './mappingEvidence.ts';
const context=(service:string,source='client',signal='logs',carrier='log')=>({name:'pwa.app.opened',service,source,signal,carrier,count:2});
test('new mapping uses the sole received service and carrier without inventing a condition',()=>{
 const mapping=initialReceivedMapping('pwa.app.opened',[context('wfs-pwa-browser')]);
 assert.equal(mapping.service,'wfs-pwa-browser');assert.ok(!('source' in mapping));assert.equal(mapping.carrier,'log');assert.equal(mapping.key,'');assert.equal(mapping.value,'');
});
test('multiple received contexts require a choice and never silently select another service or authority',()=>{
 const contexts=[context('wfs-pwa-browser'),context('worker','client'),{...context('ignored'),name:'other'},context('metrics','client','metrics','metric')];
 assert.equal(receivedMappingContexts('pwa.app.opened',contexts).length,2);
 const mapping=initialReceivedMapping('pwa.app.opened',contexts);
 assert.equal(mapping.service,'');assert.ok(!('source' in mapping));
 const next=applyReceivedContext({...mapping,step:'App opened',key:'stage.name',value:'search'},contexts[1]);
 assert.ok(!('source' in next));assert.equal(next.service,'worker');assert.equal(next.step,'App opened');assert.equal(next.key,'');assert.equal(next.value,'');
});
test('projection variants of the same context combine while absent history keeps selectors broad',()=>{
 const contexts=[context('browser'),context('browser')];
 assert.equal(receivedMappingContexts('pwa.app.opened',contexts)[0].count,4);
 assert.equal(initialReceivedMapping('pwa.app.opened',contexts).service,'browser');
 assert.equal(initialReceivedMapping('missing',contexts).service,'');
});

test('active mapping context drops legacy role selectors without losing service or carrier',()=>{
 const mapping=applyReceivedContext({event:'opened',step:'Opened',source:'browser',kind:'step',key:'',value:''},context('client-service'));
 assert.ok(!('source' in mapping));assert.equal(mapping.service,'client-service');assert.equal(mapping.carrier,'log');
});

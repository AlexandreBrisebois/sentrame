import test from 'node:test';
import assert from 'node:assert/strict';
import { acceptsSnapshot, nodeRadius, isoLocal } from './helpers.ts';
test('a response for a superseded filter cannot replace the current snapshot', () => {
  assert.equal(acceptsSnapshot(3, 4), false);
  assert.equal(acceptsSnapshot(4, 4), true);
});
test('time scales node area monotonically with a bounded cap and explicit missing marker', () => {
  assert.ok(nodeRadius(30000) > nodeRadius(1000));
  assert.equal(nodeRadius(1800001), nodeRadius(1800000));
  assert.equal(nodeRadius(null), 22);
});

test('UTC date controls preserve the stated timezone', () => { assert.equal(isoLocal('2026-10-28T12:00'), '2026-10-28T12:00:00.000Z'); });
test('summary quantiles use merged bucket counts, never averages of daily medians', async()=>{
  const {histogramQuantile}=await import('./helpers.ts');
  assert.equal(histogramQuantile({'0':9,'10':1},.5),'1.0–1.1 ms');
  assert.equal(histogramQuantile({'-1':1},.5),'<1 ms');
  assert.equal(histogramQuantile({'999':1},.5),'>24 h');
});
test('finalized period timing keeps exclusions including stages with no complete sample', async()=>{
  const {mergeSummaryGroups}=await import('./helpers.ts');
  const [g]=mergeSummaryGroups([{groups:{'d1/v1':{excluded:{'View/measured-elapsed':2},histograms:{'Vote/measured-elapsed':{'1':3}}}}},{groups:{'d1/v1':{excluded:{'View/measured-elapsed':1},histograms:{'Vote/measured-elapsed':{'1':1}}}}}]);
  assert.equal(g.excluded['View/measured-elapsed'],3);
  assert.equal(g.histograms['Vote/measured-elapsed']['1'],4);
  assert.equal(g.exclusionsAvailable,true);
});
test('displayed UTC times do not inherit the Mac timezone', async()=>{
  const {dateTime}=await import('./helpers.ts');const old=process.env.TZ;process.env.TZ='America/Toronto';
  try{const instant=Date.parse('2026-10-28T20:00:00Z');const expected=new Date(instant).toLocaleString(undefined,{dateStyle:'medium',timeStyle:'medium',timeZone:'UTC'})+' UTC';assert.equal(dateTime(instant*1e6),expected)}finally{if(old===undefined)delete process.env.TZ;else process.env.TZ=old}
});
test('discovery remains app-wide when current live evidence is empty',async()=>{
  const {discoveryChoices}=await import('./helpers.ts');const rows=discoveryChoices([{name:'vote.committed',count:10,sources:['server'],assigned:10}],[]);
  assert.equal(rows[0]?.[0],'vote.committed');assert.equal(rows[0]?.[1].count,10);assert.deepEqual([...rows[0][1].sources],['server']);
});
test('decimal evidence nanoseconds stay exact in copy while display truncates only to milliseconds',async()=>{
  const {timestampMilliseconds,evidenceCopy,dateTime}=await import('./helpers.ts');
  const event={time:'1793217600123999999',receipt:'1793217600223123456',observed:'1793217600123123456'};
  assert.equal(timestampMilliseconds(event.time),1793217600123);
  const copied=JSON.parse(evidenceCopy(event));assert.equal(copied.time,event.time);assert.equal(copied.receipt,event.receipt);assert.equal(copied.observed,event.observed);
  assert.equal(dateTime(event.time),dateTime(1793217600123));
});
test('responsive graph keeps readable physical spacing at phone and sidebar widths',async()=>{
  const {graphGeometry}=await import('./helpers.ts');
  const phone=graphGeometry(330,7),sidebar=graphGeometry(536,7);
  assert.equal(phone.width,330);assert.equal(phone.columns,2);assert.ok(phone.cellWidth>=150);assert.equal(phone.rows,4);
  assert.equal(sidebar.width,536);assert.equal(sidebar.columns,3);assert.ok(sidebar.cellWidth>=150);assert.equal(sidebar.rows,3);
  assert.equal(phone.height,phone.rows*195+50);
});
test('receipt lag leaves future source clocks unavailable and preserves positive decimal differences',async()=>{
  const {receiptLagMilliseconds}=await import('./helpers.ts');
  assert.equal(receiptLagMilliseconds('1793217600123000000','1793217600123000001'),null);
  assert.equal(receiptLagMilliseconds('1793217600123000001','1793217600123000000'),0.000001);
  assert.equal(receiptLagMilliseconds(2000000,1000000),1);
  assert.equal(receiptLagMilliseconds(Infinity,1000000),null);
});

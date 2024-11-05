import { test, expect } from '@playwright/test';

test.beforeEach(async ({ page }) => {
  await page.goto('http://localhost:3000/admin/login');
  await expect(page).toHaveURL(/login$/);
  await page.getByPlaceholder('Email address').fill('admin@seatsurfing.local');
  await page.getByRole('button', { name: '➤' }).click();
  await page.getByPlaceholder('Password').fill('12345678');
  await page.getByRole('button', { name: '➤' }).click();
  await expect(page).toHaveURL(/dashboard$/);
});

test('crud location', async ({ page }) => {
  const name = 'Location ' + Math.random().toString().substr(2);

  await page.getByRole('link', { name: 'Areas' }).click();
  await expect(page).toHaveURL(/locations$/);

  await page.getByRole('link', { name: 'Add' }).click();
  await expect(page).toHaveURL(/locations\/add$/);

  await page.getByPlaceholder('Name').fill(name);
  await page.getByPlaceholder('Description').fill(name);
  await page.locator('#check-limitConcurrentBookings').check();
  await page.getByRole('spinbutton').fill('5');
  await page.locator('input[type="file"]').setInputFiles('../server/res/floorplan.jpg');
  await page.getByRole('button', { name: 'Save' }).click();

  await page.getByRole('button', { name: 'Add space' }).click();
  await page.locator('.space-dragger').getByRole('textbox').fill('Test 1');

  await page.getByRole('button', { name: 'Add space' }).click();
  await page.locator('.space-dragger').getByRole('textbox').nth(1).fill('Test 2');

  await page.getByRole('button', { name: 'Save' }).click();
  await page.getByRole('link', { name: 'Back' }).click();
  await expect(page).toHaveURL(/locations\/.+/);
  await page.getByRole('cell', { name: name }).click();

  page.on('dialog', dialog => dialog.accept());
  await page.getByRole('button', { name: 'Delete' }).click();

  await expect(page).toHaveURL(/locations$/);
  await expect(page.getByRole('cell', { name: name })).toHaveCount(0);
});

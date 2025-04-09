import {expect, test} from '@playwright/test';

// See here how to get started:
// https://playwright.dev/docs/intro
test('visits the app root url', async ({page}) => {
  await page.goto('/');
  await expect(page.locator('h1')).toHaveText('Galvanico');
})

test('visit register page', async ({page}) => {
  await page.goto('/auth/register');
  await expect(page.locator('h2')).toHaveText('Register');
  await expect(page.locator('.btn.btn-primary.w-full')).toHaveText('Register');
})

test('visit login page', async ({page}) => {
  await page.goto('/auth/login');
  await expect(page.locator('h2')).toHaveText('Login');
  await expect(page.locator('.btn.btn-primary.w-full')).toHaveText('Login');
})

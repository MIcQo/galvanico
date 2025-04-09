<script setup lang="ts">
</script>

<template>
  <div class="hero bg-base-200 min-h-screen">
    <div class="hero-content text-center">
      <div class="max-w-md">
        <h1 class="text-5xl font-bold">Galvanico</h1>
        <p class="py-6">
          Open-source game inspired by Ikariam
        </p>
        <RouterLink :to="{name:`auth.register`}" class="btn btn-primary">Get Started</RouterLink>
      </div>
    </div>
  </div>
</template>

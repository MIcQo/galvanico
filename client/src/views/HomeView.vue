<script setup lang="ts">
import {isAuthenticated} from "@/services/auth.ts";
</script>

<template>
  <div class="hero bg-base-200 min-h-screen">
    <div class="hero-content text-center">
      <div class="max-w-md">
        <h1 class="text-5xl font-bold">Galvanico</h1>
        <p class="py-6">
          {{ $t('index.imprint') }}
        </p>
        <div>
        <RouterLink :to="{name: isAuthenticated() ? `game.index` : `auth.register`}"
                    class="btn btn-soft btn-primary">
          {{ $t('index.getStarted') }}
        </RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>

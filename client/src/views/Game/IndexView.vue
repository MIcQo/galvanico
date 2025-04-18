<script setup lang="ts">
import {computed, onMounted, onUnmounted, ref} from "vue";
import Building from "@/components/City/Building.vue";
import {useModal} from "@/stores/modal.ts";
import {useCurrentCity} from "@/stores/user.ts";
import {storeToRefs} from "pinia";

const worldview = ref<HTMLElement | null>(null);
const worldmap = ref<HTMLElement | null>(null);
const modal = useModal();

const minScale = 0.75;
const maxScale = 1;
const zoomSensitivity = 0.001; // Adjust for zoom speed
const panSensitivity = 1.5; // Adjust for pan sensitivity
const mapWidth = 5240;
const mapHeight = 1800;

let scale = 1;
let isDragging = false;
let startX = 0;
let startY = 0;
let mapLeft = 0;
let mapTop = 0;
let viewportWidth = window.innerWidth;
let viewportHeight = window.innerHeight;

const handleScroll = (e: WheelEvent): void => {
  e.preventDefault();
  if (!worldmap.value) {
    return
  }

  // Get mouse position relative to the map element
  const rect = worldmap.value.getBoundingClientRect();
  const mouseX = e.clientX - rect.left;
  const mouseY = e.clientY - rect.top;

  // Calculate the scale factor based on the wheel deltaY
  const scaleFactor = 1 - e.deltaY * zoomSensitivity;

  // Calculate the new scale, clamping it within the min and max limits
  const newScale = Math.max(minScale, Math.min(maxScale, scale * scaleFactor));

  // Calculate the change in scale
  const scaleChange = newScale / scale;

  // Update the scale
  scale = newScale;

  const currentTop = parseFloat(worldmap.value.style.top) || 0;
  const currentLeft = parseFloat(worldmap.value.style.left) || 0;

  // Calculate the vertical offset due to scaling
  const verticalOffset = (rect.height * (1 - newScale)) / 2;

  // Adjust the top position to compensate for scaling
  const adjustedTop = currentTop + verticalOffset;

  const deltaX = (mouseX - rect.width / 2) * (1 - scaleChange);
  const deltaY = (mouseY - rect.height / 2) * (1 - scaleChange);

  mapLeft = currentLeft + deltaX * panSensitivity;
  mapTop = adjustedTop - deltaY * panSensitivity;

  clampPosition()
  updateMapOffset()
};

const handleMouseDown = (e: MouseEvent): void => {
  isDragging = true;
  startX = e.clientX;
  startY = e.clientY;
}
const handleMouseUp = (): void => {
  isDragging = false;
}
const handleMouseMove = (e: MouseEvent): void => {
  if (!isDragging) return;
  if (modal.isOpen) return;

  const dx = e.clientX - startX;
  const dy = e.clientY - startY;

  startX = e.clientX;
  startY = e.clientY;

  mapLeft += dx;
  mapTop += dy;

  clampPosition()
  updateMapOffset()
}

const resizeHandler = (e: Event): void => {
  if (!worldview.value || !e.target) {
    return;
  }

  viewportWidth = window.innerWidth;
  viewportHeight = window.innerHeight;

  const w = e.target as Window;
  worldview.value.style.width = w.innerWidth + 'px';
  worldview.value.style.height = w.innerHeight + 'px';
}

const clampPosition = () => {
  const scaledWidth = mapWidth;
  const scaledHeight = mapHeight * scale;

  const minX = viewportWidth - scaledWidth;
  const minY = viewportHeight - scaledHeight;

  mapLeft = Math.min(0, Math.max(minX, mapLeft));
  mapTop = Math.min(0, Math.max(minY, mapTop));
}

const updateMapOffset = (): void => {
  if (!worldview.value) return;


  if (!worldmap.value) return;
  worldmap.value.style.top = `${mapTop}px`
  worldmap.value.style.left = `${mapLeft}px`
  worldmap.value.style.transform = `scale(${scale})`
}

onMounted(() => {
  if (worldview.value) {
    worldview.value.style.width = `${window.innerWidth}px`;
    worldview.value.style.height = `${window.innerHeight - 46}px`;
    worldview.value.addEventListener('mousedown', handleMouseDown);
  }

  if (worldmap.value) {
    // center map on mount
    mapTop = 0
    mapLeft = (viewportWidth - mapWidth) / 2;
    updateMapOffset()

    worldmap.value.addEventListener('wheel', handleScroll, {passive: false});
  }

  window.addEventListener('mouseup', handleMouseUp);
  window.addEventListener('mousemove', handleMouseMove);
  window.addEventListener("resize", resizeHandler);
})

onUnmounted(() => {
  if (worldmap.value) {
    worldmap.value.removeEventListener('wheel', handleScroll);
  }
  if (worldview.value) {
    worldview.value.removeEventListener('mousedown', handleMouseDown);
  }
  window.removeEventListener('mouseup', handleMouseUp);
  window.removeEventListener('mousemove', handleMouseMove);
  window.removeEventListener("resize", resizeHandler);
})

const {currentCity} = storeToRefs(useCurrentCity());
let emptyPositions = computed(() => {
  if (!currentCity.value) return [];
  const occupied = currentCity.value.Buildings.map(value => value.Position)
  const empty: number[] = [];

  for (let i = 0; i <= 24; i++) {
    if (occupied?.indexOf(i) !== -1) {
      continue;
    }

    empty.push(i)
  }

  return empty;
})

</script>

<template>
  <div>
    <span>{{ currentCity }}</span>
    <div ref="worldview" class="worldview">
      <div class="scollcover">
        <div ref="worldmap" class="worldmap">
          <div class="locations">
            <div class="dockyard_water hide"></div>
            <div id="city_background_nw" class="city_background"></div>
            <div id="city_background_ne" class="city_background"></div>
            <div id="city_background_sw" class="city_background"></div>
            <div id="city_background_se" class="city_background"></div>

            <Building v-for="cb of currentCity?.Buildings" :building="cb.Building.toString()"
                      :position="cb.Position"/>
            <building v-for="e in emptyPositions" building="land" :position="e"/>
          </div>
          <div class="city_water_bottom"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

.city_background {
  position: absolute;
  width: 960px;
  height: 600px;
}

#city_background_se {
  top: 600px;
  left: 960px;
  background-image: url(//gf3.geo.gfsrv.net/cdn53/a30360772bf16ff94fd68ee8907093.jpg);
}

#city_background_sw {
  top: 600px;
  left: 0;
  background-image: url(//gf1.geo.gfsrv.net/cdnc0/76268085023fbfe8b754a91dcae395.jpg);
}

#city_background_ne {
  top: 0;
  left: 960px;
  background-image: url(//gf3.geo.gfsrv.net/cdn29/180fea11f00e15aa0297c7db073cbd.jpg);
}

#city_background_nw {
  top: 0;
  left: 0;
  background-image: url(//gf3.geo.gfsrv.net/cdn5f/834d4157dbe90d8ff04f3cbab0fd7b.jpg);
}

.worldview {
  left: 0;
  padding: 0;
  position: absolute;
  top: 46px;
  z-index: 2;
  background: url(//gf1.geo.gfsrv.net/cdn3e/112fa07fd8c75ed3e2dfc28c84a9d0.jpg) repeat #E4CA99;

  .scollcover {
    height: 100%;
    overflow: hidden;
    width: 100%;
    z-index: 35;

    .worldmap {
      background: url(//gf1.geo.gfsrv.net/cdn3e/112fa07fd8c75ed3e2dfc28c84a9d0.jpg) repeat-x 0 0 transparent;
      width: 5240px;
      height: 1800px;
      cursor: move;
      position: relative;
      transition: transform 0.1s ease-out;
      transform-origin: top center;

      .locations {
        height: 1200px;
        width: 1920px;
        left: 1652px;
        margin: 0;
        overflow: hidden;
        position: absolute;
        top: 0;
        z-index: 200;
      }

      .city_water_bottom {
        background: url(//gf1.geo.gfsrv.net/cdn96/4c2935aeea148923baa9e2ba060b0b.jpg) no-repeat 0 0 transparent;
        height: 600px;
        left: 1652px;
        position: absolute;
        top: 1200px;
        width: 1920px;

        z-index: 100;
      }
    }
  }
}
</style>

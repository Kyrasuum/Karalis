<script setup lang="ts">
import { useQuasar } from 'quasar';
const q = useQuasar();
let init = 0;
if (init < 1) {
	q.dark.toggle();
	init++;
}

function toggleDarkMode() {
	q.dark.toggle();
}
</script>

<template>
	<router-view :class="[$q.dark.isActive ? 'dark': '']"/>
</template>


<script lang="ts">
import { defineComponent } from 'vue';
import { until } from 'src/types/net';

export default defineComponent({
	name: 'App',
    data() {
        return {
        };
    },
	created: function () {
	},
    methods: {
    },
});
</script>
<style scoped></style>

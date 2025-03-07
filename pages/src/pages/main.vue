<template>
	<q-layout view="hHh LpR fFf">
		<q-header elevated>
		</q-header>

		<q-page-container>
			<q-page class="" :class="[$q.dark.isActive ? 'bg-primary': '']">
			</q-page>
		</q-page-container>

		<q-footer elevated style="height: 0px; overflow: visible">
		</q-footer>
	</q-layout>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { until } from 'src/types/net';
import { url } from 'src/api/raw/raw';
import router from 'src/router';
import { useQuasar } from 'quasar';

export default defineComponent({
	name: "main",
  	components: {},
	props: {},
	setup() {},
	data() {
        const q = useQuasar();
	
		return {
			url,
			q,
		};
	},
	computed: {},
	watch: {},
	created: function () {},
	methods: {
	},
});
</script>

<style scoped lang="scss">
</style>

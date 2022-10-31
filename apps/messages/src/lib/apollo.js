import { ApolloClient, InMemoryCache } from "@apollo/client/core";
import { useSubscription } from "@vue/apollo-composable";
import gql from "graphql-tag";

const apolloCache = new InMemoryCache();

import { watch } from "vue";

// import { setContext } from "@apollo/client/link/context";
// const authLink = setContext((_, { headers }) => {
// 	return {
// 		headers: {
// 			...headers,
// 			authorization: "Foo",
// 		},
// 	};
// });

// import { HttpLink } from "@apollo/client/link/http";
// const httpLink = new HttpLink({
// 	credentials: "include",
// 	uri: import.meta.env.VITE_API_ENDPOINT,
// });

console.log(window.location);

import { WebSocketLink } from "@apollo/client/link/ws";
const wsLink = new WebSocketLink({
	uri: `wss://${window.location.hostname}/apps/messages/graph`,
	options: {
		reconnect: true,
		connectionParams: {
			credentials: "include",
		},
	},
});

const apolloClient = new ApolloClient({
	cache: apolloCache,
	credentials: "include",
	link: wsLink,
	// link: httpLink,
});

// function cacheDelete(query, field, id) {
// 	let data = apolloCache.readQuery({ query });
// 	let items = data[field].filter((i) => i.id !== id);
// 	data = { ...data };
// 	data[field] = items;
// 	apolloCache.writeQuery({ query: query, overwrite: true, data });
// }

// function cacheUpdate(query, field, item, sort = null) {
// 	let data = apolloCache.readQuery({ query });
// 	let items = [...data[field].filter((i) => i.id !== item.id), item];
// 	if (sort) items = items.sort((a, b) => a[sort].localeCompare(b[sort]));
// 	data = { ...data };
// 	data[field] = items;
// 	apolloCache.writeQuery({ query: query, overwrite: true, data });
// }

function watchType(name, refetch) {
	const { result } = useSubscription(
		gql`
			subscription ($name: String!) {
				type_changed(name: $name)
			}
		`,
		{
			name: name,
		}
	);
	watch(result, () => {
		refetch();
	});
}

export { apolloClient, apolloCache, watchType };

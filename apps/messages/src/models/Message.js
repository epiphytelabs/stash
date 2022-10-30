import { provideApolloClient } from "@vue/apollo-composable";
import { apolloClient } from "@/lib/apollo";
provideApolloClient(apolloClient);

import gql from "graphql-tag";
import { useQuery } from "@vue/apollo-composable";

export default function (hash) {
	const query = useQuery(
		gql`
			query ($hash: ID!) {
				message(hash: $hash) {
					body {
						html
						text
					}
					hash
					from {
						address
						name
					}
					received
					subject
					to
				}
			}
		`,
		{ hash }
	);

	return query;
}

import { provideApolloClient } from "@vue/apollo-composable";
import { apolloClient } from "@/lib/apollo";
provideApolloClient(apolloClient);

import gql from "graphql-tag";
import { useQuery } from "@vue/apollo-composable";
import { watchType } from "@/lib/apollo";

export default function (params) {
	const query = useQuery(
		gql`
			query {
				messages {
					from {
						address
						name
					}
					received
					subject
				}
			}
		`,
		params
	);

	watchType("message", query.refetch);

	return query;
}

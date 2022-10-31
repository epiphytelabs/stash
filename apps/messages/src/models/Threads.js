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
				threads {
					id
					messages {
						body {
							html
							text
						}
						from {
							display
						}
						received
						subject
					}
					updated
				}
			}
		`,
		params,
		{
			errorPolicy: "all",
		}
	);

	watchType("message", query.refetch);

	return query;
}

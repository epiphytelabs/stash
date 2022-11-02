import { watch } from "vue";
import { provideApolloClient } from "@vue/apollo-composable";
import { apolloClient } from "@/lib/apollo";
provideApolloClient(apolloClient);

import gql from "graphql-tag";
import { useQuery, useSubscription } from "@vue/apollo-composable";

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

	const { result: adds } = useSubscription(
		gql`
			subscription {
				threadAdded {
					id
				}
			}
		`
	);

	watch(adds, () => {
		console.log("adds", adds);
	});

	// function watchType(name, refetch) {
	// 	const { result } = useSubscription(
	// 		gql`
	// 			subscription ($name: String!) {
	// 				type_changed(name: $name)
	// 			}
	// 		`,
	// 		{
	// 			name: name,
	// 		}
	// 	);
	// 	watch(result, () => {
	// 		refetch();
	// 	});
	// }

	// watchType("message", query.refetch);

	return query;
}

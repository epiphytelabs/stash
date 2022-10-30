import { provideApolloClient } from "@vue/apollo-composable";
import { apolloClient } from "@/lib/apollo";
provideApolloClient(apolloClient);

import gql from "graphql-tag";
import { useMutation, useQuery } from "@vue/apollo-composable";

export default {
	register() {
		return useMutation(
			gql`
				mutation ($id: String!, $password: String!) {
					register(id: $id, password: $password)
				}
			`
		);
	},
	self() {
		return useQuery(
			gql`
				query {
					user {
						id
					}
				}
			`
		);
	},
};

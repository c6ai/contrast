package seedengine

import (
	"crypto/x509"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeedEngine_New(t *testing.T) {
	testCases := map[string]struct {
		secretSeed string
		salt       string

		wantPodStateSeed          string // hex encoded
		wantHistorySeed           string // hex encoded
		wantRootCAKey             string // DER, hex encoded
		wantTransactionSigningKey string // DER, hex encoded
		wantErr                   bool
	}{
		"successful 1": {
			secretSeed:       "ccebed634ddee7535cd593e1e200b19b780f3906d8782207fa09c59e87a07cb3",
			salt:             "8c1b1225c5f6cb7eef6dbd8f77a1e1e149de031d6e3718e660a8b04c8e2b0037",
			wantPodStateSeed: "e8d42dc81aea4b0d0749b75004f3bb2ad35fd827e05727ac19c31e106ddc2a1f",
			wantHistorySeed:  "e0f4adb8326ed1bbf99b8291d7a90363113e2ac8ff9d030bcabe5e48b88bf0a6",
			wantRootCAKey: "30770201010420d06cdd45c8e0a5cf51c6c7661e32304590771f24d6bdb27c81e93" +
				"d726d163464a00a06082a8648ce3d030107a1440342000498a49da1b944dccab3d8cd1409b4487" +
				"b1f2b1ceeaf15a3f42515e7cafb451368e661377b41c40c79bdb73103a47ed4b88bc0b5d2abfb5" +
				"0e6ed87446fbf3104d7",
			wantTransactionSigningKey: "30770201010420cf1baa3e9574a048c3c616add688225fe782158cb" +
				"e743e81f80938de995b3ba9a00a06082a8648ce3d030107a14403420004df88beb95e90d80dcd5" +
				"bcd5b09f95507f1e9203792b037636b12baa67fa531d380b2bf6e113783972f38743065cc2d454" +
				"7c19d0fffe35052b19db917c7761bc4",
		},
		"successful 2": {
			secretSeed:       "1adb326866d5b1e04520d9475f6ff41d3370bec96bbb5045d8dd9d16b3c48274",
			salt:             "0d2f1e0360d8476f836b8aa3dde3b1c58a361469a6e4f10cf9ed500a651c1c2b",
			wantPodStateSeed: "8964f1750fa69bd4107ffb055f3e093cd93064590f6c3ad459b66c3bb19231fd",
			wantHistorySeed:  "03c95af2f666f44239a92d2cda3a14c3ad9ad776ef06fd97a8873457b9cff7f4",
			wantRootCAKey: "30770201010420b95a0f6da02097486a042a4ba05419b747bf4ba0568b3aca303be" +
				"9bebdf700caa00a06082a8648ce3d030107a144034200046be53ed1d59f8d9c34b9ac975ce53b0" +
				"6a124b7a643ff282a0550cccdb8eb2c777a233220541a8afe88f3156496ecce428708f73950e45" +
				"7027a530df6715f6356",
			wantTransactionSigningKey: "307702010104204855f1c82608b2cd4a6350eb2865b69f16daae1c1" +
				"c7fd9337435a0ce7fc92439a00a06082a8648ce3d030107a14403420004b4a60057154c1a90a2d" +
				"721064a0db4e76e1791e2b69a1a3be40e1eed9c4584c1d1bd6fef326d17e0ab1c12e41ca1c7f4d" +
				"c3e061e7288f8110665fecee259427a",
		},
		"short salt": {
			secretSeed: "ccebed634ddee7535cd593e1e200b19b780f3906d8782207fa09c59e87a07cb3",
			salt:       "8c1b1225c5f6cb7eef6dbd8f77a1e1e149de031d6e3718e660a8b04c",
			wantErr:    true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			secret, err := hex.DecodeString(tc.secretSeed)
			require.NoError(err)
			salt, err := hex.DecodeString(tc.salt)
			require.NoError(err)

			se, err := New(secret, salt)

			if tc.wantErr {
				assert.Error(err)
				return
			}
			assert.NoError(err)
			assert.Equal(tc.wantPodStateSeed, hex.EncodeToString(se.podStateSeed))
			assert.Equal(tc.wantHistorySeed, hex.EncodeToString(se.historySeed))
			rootCAKey, err := x509.MarshalECPrivateKey(se.RootCAKey())
			require.NoError(err)
			assert.Equal(tc.wantRootCAKey, hex.EncodeToString(rootCAKey))
			transactionSigningKey, err := x509.MarshalECPrivateKey(se.TransactionSigningKey())
			require.NoError(err)
			assert.Equal(tc.wantTransactionSigningKey, hex.EncodeToString(transactionSigningKey))
		})
	}
}

func TestSeedEngine_DerivePodSecret(t *testing.T) {
	require := require.New(t)

	// policyHash -> podSecret
	testCases := map[string]string{
		"8d62644ef9944dbbb1a2b1a574840cbd6b09e5e7f96ac0f82a8a37271edd983b": "27a9ce52ad64f131d7e44c655d4ab0b0ab41b38a538615d2b28f88cbfeac2c70",
		"b838a7adb60d110d6c3c7a1dfa51b439b78386f439a092eda0d67d53cc01c02e": "257172cbb64f1681f25168d46f361aa512c08c11c21ef6ad0b7d8b46ad29d443",
		"11103d1efce19d05f5aaac2c8af405136ad91dae9f64ba25c2402100ff0e03eb": "425b229b7f327ca82ee39941cce26ea84e6a78aef3358c0c98b76515129dac32",
		"d229c5714ca84d4e73b973636723e6cd5fe49f3c3e486732facfba61f94a10fc": "9e743b32c2fb0a9d791ba4cbd51445478d118ea88c4a0953576ed1ef4c1e353f",
		"91b7513a7709d2ab92d2c1fe1e431e37f0bea18165dd908b0e6386817b0c6faf": "86343cf90cecf6a1582465d50c33a6ef38dea6ca95e1424dc0bca37d5c8e076f",
		"99704c8b2a08ae9b8165a300ea05dbeae3b4c9a2096a6221aa4175bad43d53ec": "4006cbada495cb8f95e67f1b55466d63d94ca321789090bb80f01ae6c19ce8bf",
		"f2e57529d3b92832eef960b75b2d299aadf1e373473bebf28cc105dae55c5f4e": "66d4fd6a3bfeac05490a29e6e3c4191cb2400a1949d3b4bc726a08d12415eeb5",
	}

	secretSeed, err := hex.DecodeString("ccebed634ddee7535cd593e1e200b19b780f3906d8782207fa09c59e87a07cb3")
	require.NoError(err)
	salt, err := hex.DecodeString("8c1b1225c5f6cb7eef6dbd8f77a1e1e149de031d6e3718e660a8b04c8e2b0037")
	require.NoError(err)

	se, err := New(secretSeed, salt)
	require.NoError(err)

	for policyHashStr, wantPodSecret := range testCases {
		t.Run(policyHashStr, func(t *testing.T) {
			assert := assert.New(t)

			policyHash, err := hex.DecodeString(policyHashStr)
			require.NoError(err)
			podSecret, err := se.DerivePodSecret(policyHash)
			assert.NoError(err)
			assert.Equal(wantPodSecret, hex.EncodeToString(podSecret))
		})
	}
}

func TestSeedEngine_DeriveMeshCAKey(t *testing.T) {
	req := require.New(t)

	// transactionHash -> hex(x509.MarshalECPrivateKey(meshCAKey))
	testCases := map[string]string{
		"8d62644ef9944dbbb1a2b1a574840cbd6b09e5e7f96ac0f82a8a37271edd983b": "307702010104208a9f276ddd591faec120bb350b797621963769361f256ea78cd93064c99c1b3ba00a06082a8648ce3d030107a14403420004b268c21d251bbf9dbbfb2d30e3bbc58164975e07b0853eedfdca7f96cb6013d2601ed563cb86efd5f553fe7d0fd11f923cf106e364c83f06ef0c73a9f1568842",
		"b838a7adb60d110d6c3c7a1dfa51b439b78386f439a092eda0d67d53cc01c02e": "307702010104206a708b013040b71aac3f3071681e21291b46960f5a2833610ec18afd1810c36aa00a06082a8648ce3d030107a144034200044fb57f5471a0899e79afab90a9228a4fcee338a59ab63f36b3da133c3bb269cd04bd9a4562b6d3f70c43cfae61cab95b8e5f09655299d88790ed7e5a36e86866",
		"11103d1efce19d05f5aaac2c8af405136ad91dae9f64ba25c2402100ff0e03eb": "30770201010420798e78c243df62c432c3d8a4d71a6f23d1d75a19b3699eb1b09dd287edff14f6a00a06082a8648ce3d030107a144034200047f3c507ebb83f6b126fb6064cf53912084fd8185797103c98f09218e4a2fe18245e885742ec19f468c9e1f42c43ca65d4a32881a49bfe2a478051e275c4ed307",
		"d229c5714ca84d4e73b973636723e6cd5fe49f3c3e486732facfba61f94a10fc": "3077020101042025e6cbe644f42ae6dee8d9a27adc96f7c9f5db52922b5f95c8384942860e3851a00a06082a8648ce3d030107a14403420004253cacd4362a91c3bfc9f236dd488c04a9a140792d31580f331d730db8bb03e551d5aefc04b72dc569bdd928ab4c01c42628a8da3df7e50afdb40979031e99f7",
		"91b7513a7709d2ab92d2c1fe1e431e37f0bea18165dd908b0e6386817b0c6faf": "30770201010420329f9e4a5dd04ad63dde23b61fb7991de4d2cd63c6cedd94ee318cbc10578173a00a06082a8648ce3d030107a1440342000472884cd160bedafabd874db0de1a601f00a93946c0b427cd5cbbe0c64369b40fd9ccdf79105a24e666763c78b362a1f2071e52bbce9d33c078246f7147ddeb10",
		"99704c8b2a08ae9b8165a300ea05dbeae3b4c9a2096a6221aa4175bad43d53ec": "30770201010420c51e1405af407e5c2d532e43b14c0f183a2196f09d2ab10e36dc69fd80aa230da00a06082a8648ce3d030107a14403420004303577d454826edabc474c1e7dc027f215be5749955b57bfed6374147bab7365737bb19324d78620cc0f70447b7f00ebf1a7b357b61d0424db06307b06bc2d77",
		"f2e57529d3b92832eef960b75b2d299aadf1e373473bebf28cc105dae55c5f4e": "307702010104204f50a40cea545d1b60f810fd6f517aea07f7788a8fce106e3d54dc3e3fb0fadea00a06082a8648ce3d030107a1440342000482aa46317ff7dfa4703bb61ea76bdfcd1c57d3069a7fc597f86376e48396d9a4ca7788751368b0472a9847bc3b4216efe68ed5f75a0eb158a56c69652715df6b",
	}

	secretSeed, err := hex.DecodeString("ccebed634ddee7535cd593e1e200b19b780f3906d8782207fa09c59e87a07cb3")
	req.NoError(err)
	salt, err := hex.DecodeString("8c1b1225c5f6cb7eef6dbd8f77a1e1e149de031d6e3718e660a8b04c8e2b0037")
	req.NoError(err)

	se, err := New(secretSeed, salt)
	req.NoError(err)

	for transactionHashStr, wantMeshKey := range testCases {
		t.Run(transactionHashStr, func(t *testing.T) {
			require := require.New(t)

			transactionHash, err := hex.DecodeString(transactionHashStr)
			require.NoError(err)
			meshCAKey, err := se.DeriveMeshCAKey(transactionHash)
			require.NoError(err)
			meshCAKeyDER, err := x509.MarshalECPrivateKey(meshCAKey)
			require.NoError(err)

			require.Equal(wantMeshKey, hex.EncodeToString(meshCAKeyDER))
		})
	}
}

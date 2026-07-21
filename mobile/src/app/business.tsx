import { View, Text, ScrollView, StyleSheet, TouchableOpacity, TextInput } from 'react-native';
import { useColorScheme } from 'react-native';
import { useState } from 'react';
import { Colors } from '@/constants/theme';
import { getBusinessListings, createBusinessListing } from '@/lib/api';

export default function BusinessScreen() {
  const scheme = useColorScheme();
  const colors = Colors[scheme === 'unspecified' ? 'light' : scheme];
  const [searchQuery, setSearchQuery] = useState('');
  const [showModal, setShowModal] = useState(false);

  return (
    <View style={[styles.container, { backgroundColor: colors.background }]}>
      <ScrollView style={styles.scrollView}>
        <View style={styles.header}>
          <Text style={[styles.title, { color: colors.text }]}>Business</Text>
          <Text style={[styles.subtitle, { color: colors.textSecondary }]}>
            Old Starehian businesses and services
          </Text>
        </View>

        <View style={styles.searchContainer}>
          <TextInput
            style={[styles.searchInput, { color: colors.text, backgroundColor: colors.card }]}
            placeholder="Search businesses..."
            placeholderTextColor={colors.textSecondary}
            value={searchQuery}
            onChangeText={setSearchQuery}
          />
        </View>

        <View style={styles.content}>
          <View style={styles.headerRow}>
            <Text style={[styles.sectionTitle, { color: colors.text }]}>
              Business Listings
            </Text>
            <TouchableOpacity
              style={[styles.addButton, { backgroundColor: colors.primary }]}
              onPress={() => setShowModal(true)}
            >
              <Text style={styles.addButtonText}>+ List Business</Text>
            </TouchableOpacity>
          </View>
          
          {/* Sample business cards */}
          {[
            { name: 'Starehe Tech Solutions', category: 'Technology', location: 'Nairobi', owner: 'John Doe (2015)' },
            { name: 'Elite Construction', category: 'Construction', location: 'Mombasa', owner: 'Jane Smith (2012)' },
            { name: 'Green Farms Kenya', category: 'Agriculture', location: 'Nakuru', owner: 'Peter Kamau (2018)' },
          ].map((business, index) => (
            <TouchableOpacity
              key={index}
              style={[styles.card, { backgroundColor: colors.card, borderColor: colors.border }]}
            >
              <View style={styles.cardHeader}>
                <View style={styles.businessIcon} />
                <View style={styles.cardContent}>
                  <Text style={[styles.businessName, { color: colors.text }]}>{business.name}</Text>
                  <Text style={[styles.category, { color: colors.primary }]}>{business.category}</Text>
                  <Text style={[styles.owner, { color: colors.textSecondary }]}>{business.owner}</Text>
                  <Text style={[styles.location, { color: colors.textSecondary }]}>{business.location}</Text>
                </View>
              </View>
            </TouchableOpacity>
          ))}
        </View>
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  scrollView: {
    flex: 1,
  },
  header: {
    padding: 20,
    paddingTop: 60,
  },
  title: {
    fontSize: 32,
    fontWeight: 'bold',
  },
  subtitle: {
    fontSize: 14,
    marginTop: 4,
  },
  searchContainer: {
    marginHorizontal: 20,
    marginBottom: 20,
  },
  searchInput: {
    height: 48,
    borderRadius: 12,
    paddingHorizontal: 16,
    fontSize: 14,
  },
  content: {
    padding: 20,
  },
  headerRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
  },
  addButton: {
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 8,
  },
  addButtonText: {
    color: '#fff',
    fontSize: 12,
    fontWeight: '600',
  },
  card: {
    padding: 16,
    borderRadius: 12,
    marginBottom: 12,
    borderWidth: 1,
  },
  cardHeader: {
    flexDirection: 'row',
    marginBottom: 12,
  },
  businessIcon: {
    width: 48,
    height: 48,
    borderRadius: 12,
    backgroundColor: '#4F46E5',
    marginRight: 12,
  },
  cardContent: {
    flex: 1,
  },
  businessName: {
    fontSize: 16,
    fontWeight: '600',
  },
  category: {
    fontSize: 14,
    fontWeight: '500',
    marginTop: 2,
  },
  owner: {
    fontSize: 12,
    marginTop: 4,
  },
  location: {
    fontSize: 12,
    marginTop: 2,
  },
});
